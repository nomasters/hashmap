package hashmap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// ctxKey is an unexported context key type for use with ctx
type ctxKey int

// Server Defaults
const (
	DefaultServerTimeout = 5 * time.Second
	DefaultPort          = 3000
)

// payloadCtxKey is used to give payload a collision-free key
const payloadCtxKey ctxKey = 0

var (
	s Storage
)

// ServerOptions for the hashMap Server.
type ServerOptions struct {
	Host     string
	Port     int
	TLS      bool
	CertFile string
	KeyFile  string
	Storage  StorageOptions
}

// SubmitSuccessResponse is a simple struct for json formatted submission response
type SubmitSuccessResponse struct {
	Endpoint string `json:"endpoint"`
}

// ErrorResponse a simple struct for json formatted error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// AddrString returns a string formatted as expected by the net libraries in go.
func (opts ServerOptions) AddrString() string {
	port := opts.Port
	if port == 0 {
		port = DefaultPort
	}
	return fmt.Sprintf("%v:%v", opts.Host, port)
}

// initStorage takes a StorageOptions stuct and sets the globally scoped
// Storage interface for use by the server
func initStorage(opts StorageOptions) {
	var err error
	s, err = NewStorage(opts)
	if err != nil {
		log.Fatal(err)
	}
}

// Run takes an Options struct and a server running on a specified port
// TODO: add middleware such as rate limiting and logging
func Run(opts ServerOptions) {

	// configured globally scoped storage interface
	initStorage(opts.Storage)

	r := chi.NewRouter()

	// liberal CORS support for `hashmap-client-js`
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		// ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           60,
	})

	r.Use(cors.Handler)
	r.Use(middleware.Timeout(DefaultServerTimeout))
	r.Post("/", submitHandleFunc)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"status": "healthy"}`)) })
	r.Route("/{pkHash}", func(r chi.Router) {
		r.Use(pkHashCtx)
		r.Get("/", getPayloadHandleFunc)
	})

	// run server either in http or https mode
	if opts.TLS {
		if opts.CertFile == "" || opts.KeyFile == "" {
			log.Fatal("KeyFile and CertFile are required server options to run in TLS mode")
		}
		http.ListenAndServeTLS(opts.AddrString(), opts.CertFile, opts.KeyFile, r)
	} else {
		log.Println("WARNING: running in NON-TLS MODE, this is not recommended.")
		http.ListenAndServe(opts.AddrString(), r)
	}
}

// submitHandleFunc reads and validates a Payload from r.Body and runs a series of
// Payload and Data validations. If all checks pass, the pubkey is hashed and
// the hash and payload are written to the KV store.
// TODO: return a proper JSON formatted response
func submitHandleFunc(w http.ResponseWriter, r *http.Request) {
	p, err := NewPayloadFromReader(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	d, err := p.GetData()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := d.ValidateTTL(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := d.ValidateMessageSize(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := d.ValidateTimeStamp(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	pubKey, _ := p.PubKeyBytes() // no error checking needed, already validated
	hash := MultiHashToString(pubKey)

	pwm := PayloadWithMetadata{
		Payload: *p,
		// TODO add metadata stuff
	}

	if err := s.Set(hash, pwm); err != nil {
		log.Println(err)
		http.Error(w, "internal error saving payload", 500)
		return
	}

	response, _ := json.Marshal(SubmitSuccessResponse{Endpoint: hash})

	w.Write([]byte(response))
}

// getPayloadHandleFunc gets a payload from Context, marshals the json,
// and returns the marshaled json in the response
func getPayloadHandleFunc(w http.ResponseWriter, r *http.Request) {
	// TODO: add type casting protections here
	p := r.Context().Value(payloadCtxKey).(Payload)
	payload, _ := json.Marshal(p)
	w.Write(payload)
}

// pkHashCtx is the primary response middleware used to retrieve and validate
// a payload. This middleware is designed to verify that the payload is properly
// formatted, as well as checking for common issues such as malformed hashes,
// invalid TTLs, and pubkey mismatches.
func pkHashCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pkHash := chi.URLParam(r, "pkHash")

		if err := ValidateMultiHash(pkHash); err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		pwm, err := s.Get(pkHash)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		payload := pwm.Payload

		if err := payload.Verify(); err != nil {
			// TODO: refactor to structured logs
			log.Println("payload failed to verify after reading from cache, deleting")
			log.Println(err)
			log.Println(payload)
			s.Delete(pkHash)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		// check to see if pkHash matches the actual Payload pubkey
		// this should only error if a backing store has been tampered with
		pubKey, _ := payload.PubKeyBytes() // no error checking needed, already validated
		if h := MultiHashToString(pubKey); h != pkHash {
			log.Printf("key hash does not match pubkey value hash key: %s value: %s\n", pkHash, h)
			s.Delete(pkHash)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		data, err := payload.GetData()
		if err != nil {
			log.Println("data failed to load after reading from cache, deleting")
			log.Println(err)
			log.Println(payload)
			s.Delete(pkHash)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		if err := data.ValidateTTL(); err != nil {
			log.Println(err)
			log.Println(payload)
			s.Delete(pkHash)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), payloadCtxKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

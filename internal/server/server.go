package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/nomasters/hashmap/internal/storage"
	"github.com/nomasters/hashmap/pkg/payload"
)

const (
	defaultTimeout         = 15 * time.Second
	defaultPort            = 3000
	defaultThrottleLimit   = 100
	defaultThrottleBacklog = 100
	endpointHashLength     = 88 // char count for blake2b-512 base64 string
)

// Run takes an arbitrary number of options and runs a server. The important steps here
// are that it configured a storage interface and wires that into a router to be used
// by the server handler. The runtime also leverages the Shutdown method to attempt a
// graceful shutdown in the event of an Interrupt signal. The shutdown process attempts
// to wait for all connections to close but is limited by the server timeout configuration
// which is passed into the context for the shutdown.
func Run(options ...Option) {
	var srv http.Server

	o := parseOptions(options...)
	s, err := storage.New(o.storage...)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	srv.Addr = o.addrString()
	srv.Handler = newRouter(s, o)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
		defer cancel()
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("SERVER SHUTDOWN ERROR: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Server started on: %v\n", o.addrString())
	if o.tls {
		if err := srv.ListenAndServeTLS(o.certFile, o.keyFile); err != http.ErrServerClosed {
			log.Printf("SERVER ERROR: %v", err)
		}
	} else {
		log.Println("WARNING: running in NON-TLS MODE")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("SERVER ERROR: %v", err)
		}
	}
	<-idleConnsClosed
	log.Println("\n...shutdown complete")
}

// Option is func signature used for setting Server Options
type Option func(*options)

// options contains private fields used for Option
type options struct {
	host           string
	port           int
	tls            bool
	certFile       string
	keyFile        string
	timeout        time.Duration
	limit          int
	backlog        int
	storage        []storage.Option
	allowedHeaders []string
	allowedOrigins []string
	baseRoute      string
}

// addrString returns a string formatted as expected by the net libraries in go.
func (o options) addrString() string {
	port := o.port
	if port == 0 {
		port = defaultPort
	}
	return fmt.Sprintf("%v:%v", o.host, port)
}

func newRouter(s storage.GetSetCloser, o options) http.Handler {
	r := chi.NewRouter()
	r.Use(newCors(o.allowedHeaders, o.allowedOrigins).Handler)
	r.Use(middleware.Timeout(o.timeout))
	r.Use(middleware.ThrottleBacklog(o.limit, o.backlog, o.timeout))
	r.Route(o.baseRoute, func(r chi.Router) {
		r.Use(middleware.Heartbeat("/health"))
		r.Post("/", postPayloadHandler(s))
		r.Get("/{hash}", getPayloadByHashHandler(s))
	})
	return r
}

// badRequest silently returns 400 and logs error
func badRequest(w http.ResponseWriter, v ...interface{}) {
	if len(v) > 0 {
		log.Println(v...)
	}
	http.Error(w, http.StatusText(400), http.StatusBadRequest)
}

// postPayloadHandler takes a storage.Setter and returns a http.HandlerFunc that
// uses a limited reader set to payload.MaxPayloadSize and attempts to verify
// and validate the payload in ServerMode. ServerMode verification adds an additional
// time horizon check to ensure that a payload is only written to storage within a
// strict time horizon.
func postPayloadHandler(s storage.Setter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &io.LimitedReader{R: r.Body, N: payload.MaxPayloadSize}
		body, err := ioutil.ReadAll(l)
		if err != nil {
			badRequest(w, "read error: ", err)
			return
		}
		p, err := payload.Unmarshal(body)
		if err != nil {
			badRequest(w, err)
			return
		}
		if err := p.Verify(payload.WithServerMode(true)); err != nil {
			badRequest(w, err)
			return
		}
		k := p.Endpoint()
		if err := s.Set(k, body, p.TTL, p.Timestamp); err != nil {
			badRequest(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getPayloadByHashHandler(s storage.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		k := chi.URLParam(r, "hash")
		if len(k) != endpointHashLength {
			badRequest(w, "get error. invalid hash length for:", k)
			return
		}
		if _, err := base64.URLEncoding.DecodeString(k); err != nil {
			badRequest(w, "get error. base64 decode failed for:", k, err)
			return
		}
		pb, err := s.Get(k)
		if err != nil {
			badRequest(w, "get error. storage get error for:", k, err)
			return
		}
		p, err := payload.Unmarshal(pb)
		if err != nil {
			badRequest(w, "get error. payload unmarshal failed for:", k, err)
			return
		}
		if err := p.Verify(payload.WithValidateEndpoint(k)); err != nil {
			badRequest(w, "failed get verify", k, err)
			return
		}
		w.Write(pb)
	}
}

// newCors returns cors settings with optional
func newCors(headers, origins []string) *cors.Cors {
	if len(origins) == 0 {
		origins = []string{"*"}
	}
	if len(headers) == 0 {
		headers = []string{"*"}
	}
	return cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   headers,
		AllowCredentials: true,
		MaxAge:           600,
	})
}

// parseOptions takes a arbitrary number of Option funcs and returns an options struct
func parseOptions(opts ...Option) (o options) {
	o = options{
		port:    defaultPort,
		timeout: defaultTimeout,
		limit:   defaultThrottleLimit,
		backlog: defaultThrottleBacklog,
		baseRoute: "/",
	}
	for _, option := range opts {
		option(&o)
	}

	if len(o.storage) == 0 {
		o.storage = append(o.storage, storage.WithEngine(storage.MemoryEngine))
	}

	return
}

// WithCorsAllowedHeaders takes cors pointer and returns an Option func for setting cors
func WithCorsAllowedHeaders(headers []string) Option {
	return func(o *options) {
		o.allowedHeaders = headers
	}
}

// WithCorsAllowedOrigins takes cors pointer and returns an Option func for setting cors
func WithCorsAllowedOrigins(origins []string) Option {
	return func(o *options) {
		o.allowedOrigins = origins
	}
}

// WithHost takes a string and returns an Option func for setting options.host
func WithHost(host string) Option {
	return func(o *options) {
		o.host = host
	}
}

// WithPort takes a string and returns an Option func for setting options.port
func WithPort(port int) Option {
	return func(o *options) {
		o.port = port
	}
}

// WithStorageOptions takes a set of storage options and returns an Option func for setting options.storage
func WithStorageOptions(opts ...storage.Option) Option {
	return func(o *options) {
		for _, opt := range opts {
			o.storage = append(o.storage, opt)
		}
	}
}

// WithTLS takes a boolean and returns an Option func for setting options.tls
func WithTLS(b bool) Option {
	return func(o *options) {
		o.tls = b
	}
}

// WithCertFile takes a string and returns an Option func for setting options.certFile
func WithCertFile(f string) Option {
	return func(o *options) {
		o.certFile = f
	}
}

// WithKeyFile takes a string and returns an Option func for setting options.keyFile
func WithKeyFile(f string) Option {
	return func(o *options) {
		o.keyFile = f
	}
}

// WithTimeout takes a string and returns an Option func for setting options.timeout
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}

// WithThrottle takes an int and returns an Option func for setting options.limit
func WithThrottle(t int) Option {
	return func(o *options) {
		o.limit = t
	}
}

// WithThrottleBacklog takes an int and returns an Option func for setting options.backlog
func WithThrottleBacklog(t int) Option {
	return func(o *options) {
		o.backlog = t
	}
}

// WithBaseRoute takes a string and returns an Option func for setting options.baseRoute
func WithBaseRoute(b string) Option {
	return func(o *options) {
		o.baseRoute = b
	}
}

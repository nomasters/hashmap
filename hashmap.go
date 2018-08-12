package hashmap

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/multiformats/go-multihash"
	"golang.org/x/crypto/nacl/sign"
)

const (
	MessageTTLDefault = 86400  // 1 day in seconds
	MessageTTLMax     = 604800 // 1 week in seconds
	DefaultSigMethod  = "nacl-sign-ed25519"
	Version           = "0.0.1"
	MaxPostBodySize   = 2000 // 2KB
	MaxSubmitDrift    = 15 * time.Second
	ServerTimeout     = 15 * time.Second
	DefaultPort       = ":3000"
	MaxDataBytes      = 512
	Blake2b256Code    = 45600
)

var (
	hc *HashCache
)

func init() {
	hc = NewHashCache()
}

type Payload struct {
	Message   string `json:"message"`
	Signature string `json:"sig"`
	PublicKey string `json:"pubkey"`
}

type Message struct {
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
	TTL       int64  `json:"ttl"`
	SigMethod string `json:"sigMethod"`
	Version   string `json:"version"`
}

type Validator interface {
	Validate() error
}

type NaClSignEd25519 struct {
	SignedMessage []byte
	PublicKey     *[32]byte
}

// NewNaClSignEd25519 takes to byte slices and returns a pointer to NaClSignEd25519 struct
func NewNaClSignEd25519(SignedMessage, PublicKey []byte) *NaClSignEd25519 {
	var pk [32]byte
	copy(pk[:], PublicKey)
	return &NaClSignEd25519{
		SignedMessage: SignedMessage,
		PublicKey:     &pk,
	}
}

// Validate conforms to the Validator interface and checks the validity of the NaClSignEd25519
// signed mesage against the Ed25519 pubkey
func (n NaClSignEd25519) Validate() error {
	// verify signature
	if _, valid := sign.Open(nil, append(n.SignedMessage), n.PublicKey); !valid {
		log.Printf("invalid signature: %x\n", n.PublicKey)
		return errors.New("invalid signature")
	}
	return nil
}

type Options struct {
	Port string
}

func Run(opts Options) {
	if opts.Port == "" {
		opts.Port = DefaultPort
	}
	r := chi.NewRouter()
	r.Use(middleware.Timeout(ServerTimeout))
	r.Post("/", SubmitHandleFunc)
	r.Route("/{pkHash}", func(r chi.Router) {
		r.Use(pkHashCtx)
		r.Get("/", getPayloadHandleFunc)
	})
	http.ListenAndServe(opts.Port, r)
}

// NewPayloadFromReader returns a fully verified Payload from an io.Reader source.
// This includes verifying signature and size restrictions.
func NewPayloadFromReader(r io.Reader) (*Payload, error) {
	p := Payload{}
	// read the body with strict limit on body size
	limitedReader := &io.LimitedReader{R: r, N: MaxPostBodySize}
	body, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		log.Printf("error reading payload: %v\n", err)
		return nil, errors.New("error reading payload")
	}
	// unmarshall the payload, returning an error if it fails
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("invalid payload: %v\n", err)
		return nil, errors.New("invalid payload")
	}
	if err := p.Verify(); err != nil {
		return nil, err
	}

	return &p, nil
}

// Verify validates a message signature and enforces message requirements
func (p Payload) Verify() error {
	v, err := p.NewValidator()
	if err != nil {
		return err
	}
	if err := v.Validate(); err != nil {
		return err
	}

	return nil
}

// NewValidator decodes the pubKey, Signature, and MessageBytes to byte slices.
// It then Unmarshals the Message and analyzes the SigMethod, either returning
// a Validator or an error.
func (p Payload) NewValidator() (Validator, error) {
	pubKey, err := p.PubKeyBytes()
	if err != nil {
		return nil, err
	}
	sig, err := p.SignatureBytes()
	if err != nil {
		return nil, err
	}
	messageBytes, err := p.MessageBytes()
	if err != nil {
		return nil, err
	}
	message, err := p.GetMessage()
	if err != nil {
		return nil, err
	}

	switch message.SigMethod {
	case "nacl-sign-ed25519":
		if len(pubKey) != 32 {
			return nil, errors.New("invalid pubKey length")
		}
		return NewNaClSignEd25519(append(sig, messageBytes...), pubKey), nil
	default:
		return nil, errors.New("invalid signature method")
	}

}

func (p Payload) PubKeyBytes() ([]byte, error) {
	pubKey, err := base64.StdEncoding.DecodeString(p.PublicKey)
	if err != nil {
		log.Printf("invalid pubKey encoding: %v\n", err)
		return pubKey, errors.New("invalid pubKey encoding, expecting base64")
	}
	return pubKey, nil
}

func (p Payload) SignatureBytes() ([]byte, error) {
	sig, err := base64.StdEncoding.DecodeString(p.Signature)
	if err != nil {
		log.Printf("invalid signature encoding: %v\n", err)
		return sig, errors.New("invalid signature encoding, expecting base64")
	}
	return sig, nil
}

func (p Payload) MessageBytes() ([]byte, error) {
	message, err := base64.StdEncoding.DecodeString(p.Message)
	if err != nil {
		log.Printf("invalid message encoding: %v\n", err)
		return nil, errors.New("invalid message encoding, expecting base64")
	}
	return message, nil
}

func (p Payload) GetMessage() (*Message, error) {
	// decode message
	message, err := p.MessageBytes()
	if err != nil {
		return nil, err
	}

	m := Message{}
	if err := json.Unmarshal(message, &m); err != nil {
		log.Printf("invalid message: %v\n", err)
		return nil, errors.New("invalid message")
	}

	return &m, nil
}

func (m Message) ValidateTTL() error {
	t := m.TTL

	if t > MessageTTLMax {
		return fmt.Errorf("message ttl exceeds max allowed of %v\n", MessageTTLMax)
	}

	if t == 0 {
		t = MessageTTLDefault
	}

	// convert to duration
	ttl := time.Duration(t) * time.Second

	timeStamp := time.Unix(m.Timestamp, 0)
	now := time.Now()
	diff := now.Sub(timeStamp)

	if diff > ttl {
		return errors.New("ttl exceeded")
	}

	return nil
}

func (m Message) ValidateDataSize() error {
	data, err := m.DataBytes()
	if err != nil {
		return err
	}

	if len(data) > MaxDataBytes {
		return fmt.Errorf("data exceeds max allowed: %v\n", MaxDataBytes)
	}
	return nil
}

// ValidateTimeStamp compares time.Now to message Timestamp. If the difference
// exceeds MaxSubmitDrift, it returns an error. This is to prevent replay attacks.
func (m Message) ValidateTimeStamp() error {
	timeStamp := time.Unix(m.Timestamp, 0)
	now := time.Now()
	diff := now.Sub(timeStamp)

	// get absolute value of time difference
	if diff.Seconds() < 0 {
		diff = -diff
	}

	if diff > MaxSubmitDrift {
		return errors.New("max submission time drift exceeded for message")
	}

	return nil
}

func (m Message) DataBytes() ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(m.Data)
	if err != nil {
		log.Printf("invalid data encoding: %v\n", err)
		return data, errors.New("invalid data encoding, expecting base64")
	}
	return data, nil
}

func SubmitHandleFunc(w http.ResponseWriter, r *http.Request) {
	p, err := NewPayloadFromReader(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	m, err := p.GetMessage()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := m.ValidateTTL(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := m.ValidateDataSize(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := m.ValidateTimeStamp(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	pubKey, _ := p.PubKeyBytes() // no error checking needed, already validated
	hash := MultiHashToString(pubKey)

	// USED FOR checking out values

	// log.Println("b58: " + mh.B58String())
	// log.Println("hex: " + mh.HexString())
	// dh, err := multihash.Decode(mh)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Printf("        %x\n", dh.Length)
	// log.Printf("         %x\n", dh.Code)
	// log.Printf("         %s\n", dh.Name)
	// log.Printf("             %x\n", dh.Digest)

	hc.Set(hash, *p)
	w.Write([]byte(hash))
}

// MultiHashToString takes a slice of bytes, shahes to blake2b-256
// and returns a BTC/IPFS style Base58 encoded string
func MultiHashToString(b []byte) string {
	// TODO: In the future, this should be version number aware
	mh, _ := multihash.Sum(b, Blake2b256Code, -1)
	return mh.B58String()
}

// ValidateMultiHash takes a multihash encoded in base58, decodes, and validates
// against the valid results. This may change over time, if we support more hashes.
func ValidateMultiHash(hash string) error {
	mh, err := multihash.FromB58String(hash)
	if err != nil {
		log.Printf("%v failed to decode multihash with error: \n", hash, err)
		return errors.New("multiHash Decode failed")
	}

	dh, err := multihash.Decode(mh)
	if err != nil {
		log.Printf("%v failed to decode multihash with error: \n", hash, err)
		return errors.New("multiHash Decode failed")
	}

	if dh.Length != 32 {
		return errors.New("multiHash length invalid")
	}
	if int(dh.Code) != Blake2b256Code {
		return errors.New("multiHash code invalid")
	}
	if len(dh.Digest) != 32 {
		return errors.New("pubKey hash length invalid")
	}
	return nil
}

func getPayloadHandleFunc(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value("payload").(Payload)
	payload, _ := json.Marshal(p)
	w.Write(payload)
}

func pkHashCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pkHash := chi.URLParam(r, "pkHash")

		if err := ValidateMultiHash(pkHash); err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		payload, ok := hc.Get(pkHash)
		if !ok {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		if err := payload.Verify(); err != nil {
			// TODO: refactor to structured logs
			log.Println("payload failed to verify after reading from cache, deleting")
			log.Println(err)
			log.Println(payload)
			hc.Delete(pkHash)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		// check to see if pkHash matches the actual Payload pubkey
		// this should only error if a backing store has been tampered with
		pubKey, _ := payload.PubKeyBytes() // no error checking needed, already validated
		if h := MultiHashToString(pubKey); h != pkHash {
			log.Printf("key hash does not match pubkey value hash key: %s value: %s\n", pkHash, h)
			hc.Delete(pkHash)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		message, err := payload.GetMessage()
		if err != nil {
			log.Println("message failed to load after reading from cache, deleting")
			log.Println(err)
			log.Println(payload)
			hc.Delete(pkHash)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		if err := message.ValidateTTL(); err != nil {
			log.Println(err)
			log.Println(payload)
			hc.Delete(pkHash)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), "payload", payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TODO This is the in-memory placeholder until we get a real db such as badger or redis

type HashCache struct {
	sync.RWMutex
	internal map[string]Payload
}

func NewHashCache() *HashCache {
	return &HashCache{
		internal: make(map[string]Payload),
	}
}

func (hc *HashCache) Get(key string) (Payload, bool) {
	hc.RLock()
	result, ok := hc.internal[key]
	hc.RUnlock()
	return result, ok
}

func (hc *HashCache) Set(key string, value Payload) {
	hc.Lock()
	hc.internal[key] = value
	hc.Unlock()
}

func (hc *HashCache) Delete(key string) {
	hc.Lock()
	delete(hc.internal, key)
	hc.Unlock()
}

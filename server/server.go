package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/nomasters/hashmap/storage"
)

const (
	defaultTimeout = 5 * time.Second
	defaultPort    = 3000
)

// Option is used for interacting with Context for setting Server Options
type Option func(*Context)

// Context contains private fields used for Option
type Context struct {
	host     string
	port     int
	tls      bool
	certFile string
	keyFile  string
	timeout  time.Duration
	storage  []storage.Option
}

// addrString returns a string formatted as expected by the net libraries in go.
func (c Context) addrString() string {
	port := c.port
	if port == 0 {
		port = defaultPort
	}
	return fmt.Sprintf("%v:%v", c.host, port)
}

func newRouter(s storage.Storage) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second)
		w.Write([]byte(`{"status": "healthy"}`))
	})
	return r
}

// parseOptions takes a arbitrary number of Option funcs and returns a Context
func parseOptions(options ...Option) Context {
	c := Context{
		port:    defaultPort,
		timeout: defaultTimeout,
	}
	for _, option := range options {
		option(&c)
	}

	if len(c.storage) == 0 {
		c.storage = append(c.storage, storage.WithEngine(storage.MemoryEngine))
	}

	return c
}

// WithHost takes a string and returns an Option func for setting Context.host
func WithHost(h string) Option {
	return func(c *Context) {
		c.host = h
	}
}

// WithPort takes a string and returns an Option func for setting Context.port
func WithPort(p int) Option {
	return func(c *Context) {
		c.port = p
	}
}

// WithStorageOptions takes a set of storage options and returns an Option func for setting Context.storage
func WithStorageOptions(options ...storage.Option) Option {
	return func(c *Context) {
		for _, o := range options {
			c.storage = append(c.storage, o)
		}
	}
}

// WithTLS takes a boolean and returns an Option func for setting Context.tls
func WithTLS(b bool) Option {
	return func(c *Context) {
		c.tls = b
	}
}

// WithCertFile takes a string and returns an Option func for setting Context.certFile
func WithCertFile(f string) Option {
	return func(c *Context) {
		c.certFile = f
	}
}

// WithKeyFile takes a string and returns an Option func for setting Context.keyFile
func WithKeyFile(f string) Option {
	return func(c *Context) {
		c.keyFile = f
	}
}

// WithTimeout takes a string and returns an Option func for setting Context.WithTimeout
func WithTimeout(d time.Duration) Option {
	return func(c *Context) {
		c.timeout = d
	}
}

// Run takes an arbitrary number of options and runs a server
func Run(options ...Option) {
	var srv http.Server

	c := parseOptions(options...)
	s, err := storage.NewStorage(c.storage...)
	if err != nil {
		log.Fatal(err)
	}

	srv.Addr = c.addrString()
	srv.Handler = newRouter(s)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if c.tls {
		log.Printf("Server started on: %v\n", c.addrString())
		if err := srv.ListenAndServeTLS(c.certFile, c.keyFile); err != http.ErrServerClosed {
			log.Printf("SERVER ERROR: %v", err)
		}
	} else {
		log.Println("WARNING: running in NON-TLS MODE")
		log.Printf("Server started on: %v\n", c.addrString())
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("SERVER ERROR: %v", err)
		}
	}
	<-idleConnsClosed
	log.Println("\n...shutdown complete")
}

package jnsd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	// ErrNotRegistered is a generic error to be returned if a name or
	// address is not registered.
	ErrNotRegistered = errors.New("not registered")
)

// HandlerConfig defines necessary configuration for running the
// nameserver handler.
type HandlerConfig struct {
	// Name is a function which maps from a name to an address. It must
	// not be nil.
	Name func(name string) (string, error)

	// Addr is a function which maps from an address to a name. It must
	// not be nil.
	Addr func(addr string) (string, error)
}

// Handler returns an http.Handler from a config.
func (config HandlerConfig) Handler() http.Handler {
	router := mux.NewRouter()

	router.
		Methods("GET", "OPTIONS").
		Path("/name/{name}").
		Handler(handleGet(func(n, a string) (name, addr string, err error) {
			log.Printf("request to resolve name: %q", n)
			addr, err = config.Name(n)
			return n, addr, err
		}))

	router.
		Methods("GET", "OPTIONS").
		Path("/addr/{addr}").
		Handler(handleGet(func(n, a string) (name, addr string, err error) {
			log.Printf("request to resolve address: %q", a)
			name, err = config.Addr(a)
			return name, "", err
		}))

	router.Use(mux.CORSMethodMiddleware(router))

	return router
}

// handleGet returns a handler that handles the nameserver's GET
// endpoints.
func handleGet(f func(name, addr string) (string, string, error)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(rw)

		vars := mux.Vars(req)
		name, addr, err := f(vars["name"], vars["addr"])
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			enc.Encode(Error{Error: err.Error()})
			return
		}
		if addr != "" {
			addr = fmt.Sprintf("0x%v", addr)
		}

		enc.Encode(Response{
			Name: name,
			Addr: addr,
		})
	})
}

// Response is the schema for the data sent as a successful response
// to a request to one of the endpoints.
type Response struct {
	Name string `json:"name,omitempty"`
	Addr string `json:"addr,omitempty"`
}

// Error is the schema for the data sent as an error response to a
// request to one of the endpoints.
type Error struct {
	Error string `json:"error"`
}

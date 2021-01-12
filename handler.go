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
	// NotRegistered is a generic error to be returned if a name or
	// address is not registered.
	NotRegistered = errors.New("not registered")
)

type HandlerConfig struct {
	Name func(name string) (string, error)
	Addr func(addr string) (string, error)
}

func (config HandlerConfig) Handler() http.Handler {
	router := mux.NewRouter()

	router.
		Methods("GET", "OPTIONS").
		Path("/name/{name}").
		Handler(get(func(n, a string) (name, addr string, err error) {
			log.Printf("request to resolve name: %q", n)
			addr, err = config.Name(n)
			return n, addr, err
		}))

	router.
		Methods("GET", "OPTIONS").
		Path("/addr/{addr}").
		Handler(get(func(n, a string) (name, addr string, err error) {
			log.Printf("request to resolve address: %q", a)
			name, err = config.Addr(a)
			return name, "", err
		}))

	router.Use(mux.CORSMethodMiddleware(router))

	return router
}

func get(f func(name, addr string) (string, string, error)) http.Handler {
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

type Response struct {
	Name string `json:"name,omitempty"`
	Addr string `json:"addr,omitempty"`
}

type Error struct {
	Error string `json:"error"`
}

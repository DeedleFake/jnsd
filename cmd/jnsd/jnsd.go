package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/DeedleFake/jnsd"
	"github.com/DeedleFake/jnsd/internal/cli"
)

type nameMapping struct {
	Name map[string]string
	Addr map[string]string
}

func loadMapping(path string) (nameMapping, error) {
	file, err := os.Open(path)
	if err != nil {
		return nameMapping{}, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nameMapping{}, fmt.Errorf("read: %w", err)
	}

	var mapping nameMapping
	err = json.Unmarshal(buf, &mapping.Name)
	if err != nil {
		return nameMapping{}, fmt.Errorf("unmarshal: %w", err)
	}
	mapping.updateAddr()

	return mapping, nil
}

func (mapping *nameMapping) updateAddr() {
	mapping.Addr = make(map[string]string, len(mapping.Name))
	for name, addr := range mapping.Name {
		mapping.Addr[addr] = name
	}
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	addr := flag.String("addr", ":8080", "address to listen on")
	tlscert := flag.String("tls.cert", "", "TLS certificate")
	tlskey := flag.String("tls.key", "", "TLS key")
	names := flag.String("names", "", "path to JSON name-to-address mapping")
	flag.Parse()

	listenAndServe := (*http.Server).ListenAndServe
	if (*tlscert != "") && (*tlskey != "") {
		listenAndServe = func(server *http.Server) error {
			return server.ListenAndServeTLS(*tlscert, *tlskey)
		}
	}

	mapping, err := loadMapping(*names)
	if err != nil {
		return fmt.Errorf("load names: %w", err)
	}

	handler := jnsd.HandlerConfig{
		Name: func(name string) (string, error) {
			addr, ok := mapping.Name[name]
			if !ok {
				return "", jnsd.NotRegistered
			}
			return addr, nil
		},

		Addr: func(addr string) (string, error) {
			name, ok := mapping.Addr[addr]
			if !ok {
				return "", jnsd.NotRegistered
			}
			return name, nil
		},
	}.Handler()

	server := &http.Server{
		Addr:    *addr,
		Handler: handler,
		BaseContext: func(lis net.Listener) context.Context {
			return ctx
		},
	}
	go func() {
		<-ctx.Done()

		ctx, _ := context.WithTimeout(context.Background(), time.Minute)
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("Error: shutdown server: %v", err)
		}
	}()

	log.Printf("Serving on %v", *addr)
	defer log.Printf("Shutting down...")

	err = listenAndServe(server)
	if err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func main() {
	ctx, cancel := cli.SignalContext(context.Background(), os.Interrupt)
	defer cancel()

	err := run(ctx)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

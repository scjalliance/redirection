package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/gentlemanautomaton/signaler"
	"github.com/scjalliance/redirection"
)

func main() {
	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)
	defer shutdown.Wait()
	defer shutdown.Trigger()

	cfg := DefaultConfig
	if err := cfg.ParseEnv(); err != nil {
		log.Fatalf("redirector: unable to parse environment: %v", err)
	}
	cfg.ParseArgs(os.Args[1:], flag.ExitOnError)

	if cfg.HTTPAddr == "" && cfg.HTTPSAddr == "" {
		log.Fatalf("redirector: no interfaces defined: both http and https addresses are empty")
	}

	mapper, err := load(cfg.DataFile)
	if err != nil {
		log.Fatal(err)
	}

	ctx := shutdown.Context()

	var wg sync.WaitGroup
	errchan := make(chan error)
	defer close(errchan)

	// HTTP Server
	if cfg.HTTPAddr != "" {
		wg.Add(1)
		go run(ctx, &wg, cfg.HTTPAddr, false, mapper, cfg.CertificateCacheDir, cfg.ShutdownTimeout, errchan)
	}

	// HTTPS Server
	if cfg.HTTPSAddr != "" {
		wg.Add(1)
		go run(ctx, &wg, cfg.HTTPSAddr, true, mapper, cfg.CertificateCacheDir, cfg.ShutdownTimeout, errchan)
	}

	go func() {
		for err := range errchan {
			log.Println(err)
			shutdown.Trigger()
		}
	}()

	wg.Wait()
}

func run(ctx context.Context, wg *sync.WaitGroup, addr string, useTLS bool, mapper redirection.Mapper, certCacheDir string, shutdownTimeout time.Duration, errchan chan<- error) {
	defer wg.Done()

	if ctxDone(ctx) {
		return
	}

	var (
		handler = redirection.NewHandler(mapper)
		server  = &http.Server{
			Handler: handler,
		}
	)

	if useTLS {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: redirection.HostWhitelist(mapper),
		}
		if certCacheDir != "" {
			m.Cache = autocert.DirCache(certCacheDir)
		}
		server.TLSConfig = m.TLSConfig()
	}

	if useTLS {
		log.Printf("redirector: server starting on %s with TLS", addr)
	} else {
		log.Printf("redirector: server starting on %s", addr)
	}

	err := runServer(ctx, server, useTLS, shutdownTimeout)
	errchan <- fmt.Errorf("redirector: server on %s exited: %s", addr, err)
}

func runServer(ctx context.Context, s *http.Server, useTLS bool, shutdownTimeout time.Duration) error {
	if ctxDone(ctx) {
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go shutdownServer(ctx, s, shutdownTimeout) // Shutdown the server if ctx is cancelled

	if useTLS {
		return s.ListenAndServeTLS("", "")
	}
	return s.ListenAndServe()
}

func shutdownServer(ctx context.Context, s *http.Server, timeout time.Duration) {
	// Wait a little bit to make sure the server has spun up
	time.Sleep(time.Millisecond * 100)

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.Shutdown(shutdownCtx)
}

func load(path string) (mapper redirection.Mapper, err error) {
	contents, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		return nil, fmt.Errorf("redirector: unable to read data file \"%s\": %v", path, fileErr)
	}

	// TODO: Use json.Decoder and stream the file instead of slurping?

	var element redirection.Element
	dataErr := json.Unmarshal(contents, &element)
	if dataErr != nil {
		return nil, fmt.Errorf("redirector: decoding error while parsing data file \"%s\": %v", path, dataErr)
	}

	return &element, nil
}

func ctxDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

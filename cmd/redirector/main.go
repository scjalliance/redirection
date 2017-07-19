package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/scjalliance/redirection"
)

const (
	// DefaultShutdownTimeout is the default amount of time that a graceful
	// server shutdown will be allowed
	DefaultShutdownTimeout = time.Second
)

var (
	config          = "redirect.json"
	httpAddr        = ":80"
	httpsAddr       = ":443"
	cache           = "redirector-cache"
	shutdownTimeout = DefaultShutdownTimeout
)

func main() {
	flag.StringVar(&config, "config", "redirect.json", "configuration file")
	flag.StringVar(&httpAddr, "http", httpAddr, "HTTP server listening address (blank for none)")
	flag.StringVar(&httpsAddr, "https", httpsAddr, "HTTPS server listening address (blank for none)")
	//flag.BoolVar(&acme, "acme", acme, "use ACME for automated certificate renewal")
	flag.Parse()

	if httpAddr == "" && httpsAddr == "" {
		log.Fatal("redirector: no interfaces defined: both http and https addresses are empty")
	}

	mapper, err := load(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()
	go func() {
		waitForSignal()
		shutdown()
	}()

	var wg sync.WaitGroup
	errchan := make(chan error)
	defer close(errchan)

	// HTTP Server
	if httpAddr != "" {
		wg.Add(1)
		go run(ctx, &wg, httpAddr, false, mapper, errchan)
	}

	// HTTPS Server
	if httpsAddr != "" {
		wg.Add(1)
		go run(ctx, &wg, httpsAddr, true, mapper, errchan)
	}

	go func() {
		for err := range errchan {
			log.Println(err)
		}
	}()

	wg.Wait()
}

func run(ctx context.Context, wg *sync.WaitGroup, addr string, useTLS bool, mapper redirection.Mapper, errchan chan<- error) {
	defer wg.Done()

	handler := redirection.NewHandler(mapper)

	for {
		if ctxDone(ctx) {
			return
		}

		s := &http.Server{
			Handler: handler,
		}

		if useTLS {
			m := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: redirection.HostWhitelist(mapper),
			}
			s.TLSConfig = &tls.Config{
				GetCertificate: m.GetCertificate,
				//NextProtos:     []string{"h2", "http/1.1"}, // Enable HTTP/2
			}
		}

		if useTLS {
			log.Printf("redirector: server starting on %s with TLS", addr)
		} else {
			log.Printf("redirector: server starting on %s", addr)
		}

		err := runServer(ctx, s, useTLS)
		errchan <- fmt.Errorf("redirector: server on %s exited: %s", addr, err)
		time.Sleep(time.Second * 1)
	}
}

func runServer(ctx context.Context, s *http.Server, useTLS bool) error {
	if ctxDone(ctx) {
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go shutdownServer(ctx, s) // Shutdown the server if ctx is cancelled

	if useTLS {
		return s.ListenAndServeTLS("", "")
	}
	return s.ListenAndServe()
}

func shutdownServer(ctx context.Context, s *http.Server) {
	// Wait a little bit to make sure the server has spun up
	time.Sleep(time.Millisecond * 100)

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	s.Shutdown(shutdownCtx)
}

func load(path string) (mapper redirection.Mapper, err error) {
	contents, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		return nil, fmt.Errorf("redirector: unable to read configuration file \"%s\": %v", path, fileErr)
	}

	// TODO: Use json.Decoder and stream the file instead of slurping?

	var element redirection.Element
	dataErr := json.Unmarshal(contents, &element)
	if dataErr != nil {
		return nil, fmt.Errorf("redirector: decoding error while parsing configuration file \"%s\": %v", path, dataErr)
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

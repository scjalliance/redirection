package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gentlemanautomaton/bindflag"
)

const (
	// DefaultDataFile is the default redirection map data file.
	DefaultDataFile = "redirect.json"
	// DefaultCertificateCacheDir is the default certificate cache directory.
	DefaultCertificateCacheDir = "cert-cache"
	// DefaultHTTPAddr is the default HTTP listener address.
	DefaultHTTPAddr = ":80"
	// DefaultHTTPSAddr is the default HTTP over TLS listener address.
	DefaultHTTPSAddr = ":443"
	// DefaultShutdownTimeout is the default amount of time that a graceful
	// server shutdown will be allowed.
	DefaultShutdownTimeout = time.Second
)

// Config holds a set of redirector configuration values.
type Config struct {
	DataFile            string
	CertificateCacheDir string
	HTTPAddr            string
	HTTPSAddr           string
	ShutdownTimeout     time.Duration
}

// DefaultConfig holds the default configuration values.
var DefaultConfig = Config{
	DataFile:            DefaultDataFile,
	CertificateCacheDir: DefaultCertificateCacheDir,
	HTTPAddr:            DefaultHTTPAddr,
	HTTPSAddr:           DefaultHTTPSAddr,
	ShutdownTimeout:     DefaultShutdownTimeout,
}

// ParseEnv will parse environment variables and apply them to the
// configuration.
func (c *Config) ParseEnv() error {
	var (
		data, hasData       = os.LookupEnv("DATA_FILE")
		cache, hasCache     = os.LookupEnv("CERTIFICATE_CACHE_DIR")
		httpAddr, hasHTTP   = os.LookupEnv("HTTP_ADDR")
		httpsAddr, hasHTTPS = os.LookupEnv("HTTPS_ADDR")
		sdt, hasSDT         = os.LookupEnv("SHUTDOWN_TIMEOUT")
	)
	if hasData {
		c.DataFile = data
	}
	if hasCache {
		c.CertificateCacheDir = cache
	}
	if hasHTTP {
		c.HTTPAddr = httpAddr
	}
	if hasHTTPS {
		c.HTTPSAddr = httpsAddr
	}
	if hasSDT {
		var err error
		c.ShutdownTimeout, err = time.ParseDuration(sdt)
		if err != nil {
			return fmt.Errorf("invalid value \"%s\" for variable %s: %v", sdt, "SHUTDOWN_TIMEOUT", err)
		}
	}
	return nil
}

// ParseArgs parses the given argument list and applies them to the
// configuration.
func (c *Config) ParseArgs(args []string, errorHandling flag.ErrorHandling) error {
	fs := flag.NewFlagSet("", errorHandling)
	c.Bind(fs)
	return fs.Parse(args)
}

// Bind will bind the given flag set to the configuration.
func (c *Config) Bind(fs *flag.FlagSet) {
	fs.Var(bindflag.String(&c.DataFile), "data", "redirection data file")
	fs.Var(bindflag.String(&c.CertificateCacheDir), "cache", "certificate cache directory for storage of ACME certificates")
	fs.Var(bindflag.String(&c.HTTPAddr), "http", "http server listening address")
	fs.Var(bindflag.String(&c.HTTPSAddr), "https", "https server listening address")
	fs.Var(bindflag.Duration(&c.ShutdownTimeout), "shutdownTimeout", "maximum time allowed for graceful server shutdown")
}

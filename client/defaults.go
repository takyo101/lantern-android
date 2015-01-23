package client

import (
	"crypto/x509"
	"log"

	"github.com/getlantern/fronted"
)

const (
	cloudflare = `cloudflare`
)

const (
	defaultInsecureSkipVerify = false
	defaultBufferRequest      = false
	defaultDialTimeoutMillis  = 0
	defaultRedialAttempts     = 2
	defaultWeight             = 1000000
	defaultQOS                = 10
)

var defaultCertPool *x509.CertPool

// defaultFrontedServerList holds the list of fronted servers.
var defaultFrontedServerList = []frontedServer{
	frontedServer{
		Host:          "roundrobin.getiantem.org",
		Port:          443,
		MasqueradeSet: cloudflare,
		QOS:           10,
		Weight:        1000000,
	},
}

// masqueradeSets holds a map of masquerades for fronted servers.
var masqueradeSets = map[string][]*fronted.Masquerade{
	// See masquerades.go
	cloudflare: cloudflareMasquerades,
}

func init() {
	// Populating defaultCertPool.
	var err error
	if defaultCertPool, err = getTrustedCertPool(); err != nil {
		log.Printf("getTrustedCertPool: %q", err)
	}
}

package client

import (
	"crypto/x509"
	"github.com/getlantern/keyman"
)

// CA represents a certificate authority
type CA struct {
	CommonName string
	Cert       string // PEM-encoded
}

// getTrustedCerts returns a slice of PEM-encoded certs for the trusted CAs.
func getTrustedCerts() []string {
	certs := make([]string, 0, len(defaultTrustedCAs))

	for _, ca := range defaultTrustedCAs {
		certs = append(certs, ca.Cert)
	}

	return certs
}

// getTrustedCertPool returns a certificate pool containing the trusted CAs.
func getTrustedCertPool() (certPool *x509.CertPool, err error) {
	trustedCerts := getTrustedCerts()

	if certPool, err = keyman.PoolContainingCerts(trustedCerts...); err != nil {
		return nil, err
	}

	return certPool, nil
}

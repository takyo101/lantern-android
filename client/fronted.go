package client

import (
	"fmt"
	"log"

	"github.com/getlantern/balancer"
	"github.com/getlantern/fronted"
)

type frontedServer struct {
	Host          string
	Port          int
	MasqueradeSet string
	QOS           int
	Weight        int
}

// Wraps a fronted.Dialer with a balancer.Dialer.
func (s *frontedServer) dialer() *balancer.Dialer {
	fd := fronted.NewDialer(&fronted.Config{
		Host:               s.Host,
		Port:               s.Port,
		Masquerades:        masqueradeSets[s.MasqueradeSet],
		InsecureSkipVerify: defaultInsecureSkipVerify,
		BufferRequests:     defaultBufferRequest,
		DialTimeoutMillis:  defaultDialTimeoutMillis,
		RedialAttempts:     defaultRedialAttempts,
		RootCAs:            defaultCertPool,
	})

	masqueradeQualifier := ""

	if s.MasqueradeSet != "" {
		masqueradeQualifier = fmt.Sprintf(" using masquerade set %s", s.MasqueradeSet)
	}

	return &balancer.Dialer{
		Label:  fmt.Sprintf("fronted proxy at %s:%d%s", s.Host, s.Port, masqueradeQualifier),
		Weight: s.Weight,
		QOS:    s.QOS,
		Dial:   fd.Dial,
		OnClose: func() {
			err := fd.Close()
			if err != nil {
				log.Printf("Unable to close fronted dialer: %s", err)
			}
		},
	}

}

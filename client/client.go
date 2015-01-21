package client

import (
	"fmt"
	"github.com/getlantern/balancer"
	"github.com/getlantern/fronted"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type frontedServer struct {
	Host string
	Port int
}

func (s *frontedServer) dialer() *balancer.Dialer {
	fd := fronted.NewDialer(&fronted.Config{
		Host: s.Host,
		Port: s.Port,
	})
	masqueradeQualifier := ""
	return &balancer.Dialer{
		Label:  fmt.Sprintf("fronted proxy at %s:%d%s", s.Host, s.Port, masqueradeQualifier),
		Weight: 1,
		QOS:    0,
		Dial:   fd.Dial,
		OnClose: func() {
			err := fd.Close()
			if err != nil {
				log.Printf("Unable to close fronted dialer: %s", err)
			}
		},
	}
}

// Client is a HTTP proxy that accepts connections from local programs and
// proxies these via remote flashlight servers.
type Client struct {
	Addr           string
	frontedServers []*frontedServer
	ln             *Listener
	bal            *balancer.Balancer
}

func NewClient(addr string) *Client {
	client := &Client{Addr: addr}

	client.frontedServers = make([]*frontedServer, 0, 8)

	// TODO: How are we going to add more than one fronted servers?
	client.frontedServers = append(client.frontedServers, &frontedServer{
		Host: "roundrobin.getiantem.org",
		Port: 443,
	})

	client.bal = client.initBalancer()

	return client
}

func (client *Client) getBalancer() *balancer.Balancer {
	// TODO
	return client.bal
}

func (client *Client) initBalancer() *balancer.Balancer {
	dialers := make([]*balancer.Dialer, 0, len(client.frontedServers))

	for _, s := range client.frontedServers {
		dialer := s.dialer()
		dialers = append(dialers, dialer)
	}

	bal := balancer.New(dialers...)

	return bal
}

func (client *Client) getReverseProxy() *httputil.ReverseProxy {
	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// do nothing
		},
		Transport: &http.Transport{
			// We disable keepalives because some servers pretend to support
			// keep-alives but close their connections immediately, which
			// causes an error inside ReverseProxy.  This is not an issue
			// for HTTPS because  the browser is responsible for handling
			// the problem, which browsers like Chrome and Firefox already
			// know to do.
			//
			// See https://code.google.com/p/go/issues/detail?id=4677
			DisableKeepAlives: true,
			// TODO: would be good to make this sensitive to QOS, which
			// right now is only respected for HTTPS connections. The
			// challenge is that ReverseProxy reuses connections for
			// different requests, so we might have to configure different
			// ReverseProxies for different QOS's or something like that.
			Dial: client.bal.Dial,
		},
		// Set a FlushInterval to prevent overly aggressive buffering of
		// responses, which helps keep memory usage down
		FlushInterval: 250 * time.Millisecond,
	}

	return rp
}

// ServeHTTP implements the method from interface http.Handler using the latest
// handler available from getHandler() and latest ReverseProxy available from
// getReverseProxy().
func (client *Client) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "CONNECT" {
		client.intercept(resp, req)
	} else {
		client.getReverseProxy().ServeHTTP(resp, req)
	}
}

// ListenAndServe spawns the HTTP proxy and makes it listen for incoming
// connections.
func (c *Client) ListenAndServe() (err error) {
	addr := c.Addr

	if addr == "" {
		addr = ":http"
	}

	if c.ln, err = NewListener(addr); err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    c.Addr,
		Handler: c,
	}

	log.Printf("Starting proxy server at %s...", addr)

	return httpServer.Serve(c.ln)
}

func targetQOS(req *http.Request) int {
	return 0
}

// intercept intercepts an HTTP CONNECT request, hijacks the underlying client
// connetion and starts piping the data over a new net.Conn obtained from the
// given dial function.
func (client *Client) intercept(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		panic("Intercept used for non-CONNECT request!")
	}

	// Hijack underlying connection
	clientConn, _, err := resp.(http.Hijacker).Hijack()
	if err != nil {
		respondBadGateway(resp, fmt.Sprintf("Unable to hijack connection: %s", err))
		return
	}
	defer clientConn.Close()

	addr := hostIncludingPort(req, 443)

	// Establish outbound connection
	connOut, err := client.getBalancer().DialQOS("tcp", addr, targetQOS(req))
	if err != nil {
		respondBadGateway(clientConn, fmt.Sprintf("Unable to handle CONNECT request: %s", err))
		return
	}
	defer connOut.Close()

	// Pipe data
	pipeData(clientConn, connOut, req)
}

// Stop is currently not implemented but should make the listener stop
// accepting new connections and then kill all active connections.
func (c *Client) Stop() error {
	log.Printf("Stopping proxy server...")
	return nil
}

func respondBadGateway(w io.Writer, msg string) error {
	log.Printf("Responding BadGateway: %v", msg)
	resp := &http.Response{
		StatusCode: 502,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	err := resp.Write(w)
	if err == nil {
		_, err = w.Write([]byte(msg))
	}
	return err
}

// hostIncludingPort extracts the host:port from a request.  It fills in a
// a default port if none was found in the request.
func hostIncludingPort(req *http.Request, defaultPort int) string {
	_, port, err := net.SplitHostPort(req.Host)
	if port == "" || err != nil {
		return req.Host + ":" + strconv.Itoa(defaultPort)
	} else {
		return req.Host
	}
}

// pipeData pipes data between the client and proxy connections.  It's also
// responsible for responding to the initial CONNECT request with a 200 OK.
func pipeData(clientConn net.Conn, connOut net.Conn, req *http.Request) {
	// Start piping to proxy
	go io.Copy(connOut, clientConn)

	// Respond OK
	err := respondOK(clientConn, req)
	if err != nil {
		log.Printf("Unable to respond OK: %s", err)
		return
	}

	// Then start coyping from out to client
	io.Copy(clientConn, connOut)
}

func respondOK(writer io.Writer, req *http.Request) error {
	defer req.Body.Close()
	resp := &http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	return resp.Write(writer)
}

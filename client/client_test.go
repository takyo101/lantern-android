package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const listenProxyAddr = "127.0.0.1:9997"

const expectedBody = "Google is built by a large team of engineers, designers, researchers, robots, and others in many different sites across the globe. It is updated continuously, and built with more tools and technologies than we can shake a stick at. If you'd like to help us out, see google.com/careers.\n"

func testReverseProxy() error {
	var req *http.Request

	req = &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "www.google.com",
			Path:   "http://www.google.com/humans.txt",
		},
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Host": {"www.google.com:80"},
		},
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(n, a string) (net.Conn, error) {
				//return net.Dial("tcp", "127.0.0.1:9898")
				return net.Dial("tcp", listenProxyAddr)
			},
		},
	}

	var res *http.Response
	var err error

	if res, err = client.Do(req); err != nil {
		return err
	}

	var buf []byte

	buf, err = ioutil.ReadAll(res.Body)

	fmt.Printf(string(buf))

	if string(buf) != expectedBody {
		return errors.New("Expecting another response.")
	}

	return nil
}

func TestListenAndServeStop(t *testing.T) {

	c := NewClient(listenProxyAddr)

	go func() {
		c.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)

	c.Stop()
}

func TestListenAndServeSpawn(t *testing.T) {

	go func() {
		c := NewClient(":9997")
		var err error
		if err = c.ListenAndServe(); err != nil {
			t.Fatal(err)
		}
	}()

}

func TestListenAndServeProxy(t *testing.T) {
	err := testReverseProxy()
	if err != nil {
		t.Fatal(err)
	}
}

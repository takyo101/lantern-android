package client

import (
	"compress/gzip"
	"errors"
	"github.com/getlantern/fronted"
	"github.com/getlantern/yaml"
	"io"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Client struct {
		FrontedServers []frontedServer                  `yaml:"frontedservers"`
		MasqueradeSets map[string][]*fronted.Masquerade `yaml:"masqueradesets"`
	} `yaml:"client"`
	TrustedCAs []*CA `yaml:"trustedcas"`
}

var (
	ErrFailedConfigRequest = errors.New(`Could not get configuration file.`)
)

const (
	remoteConfigURL = `https://s3.amazonaws.com/lantern_config/cloud.1.6.0.yaml.gz`
)

func pullConfigFile() ([]byte, error) {
	var err error
	var res *http.Response

	// Issuing a post request to download configuration file.
	if res, err = http.Get(remoteConfigURL); err != nil {
		return nil, err
	}

	// Expecting 200 OK
	if res.StatusCode != http.StatusOK {
		return nil, ErrFailedConfigRequest
	}

	// Using a gzip reader as we're getting a compressed file.
	var body io.ReadCloser
	if body, err = gzip.NewReader(res.Body); err != nil {
		return nil, err
	}

	// Returning uncompressed bytes.
	return ioutil.ReadAll(body)
}

func pullConfig() (*Config, error) {
	var err error
	var buf []byte

	var cfg Config

	if buf, err = pullConfigFile(); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

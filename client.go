package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	dphttp "github.com/ONSdigital/dp-net/http"
	"io/ioutil"
	"net/http"
	"strings"
)

// Nomad represents the nomad client
type Nomad struct {
	Client dphttp.Clienter
	URL string
	// Name needed for checker later.
	// Name string
}

const prefix = "https://"

var caCertPool *x509.CertPool

// NewClient returns a Nomad HTTP client for this endpoint
// with optional TLS config
func NewClient(nomadEndpoint, nomadCACert string, nomadTLSSkipVerify bool) (*Nomad, error){

	dpHTTPClient := *dphttp.DefaultClient

	if strings.HasPrefix(nomadEndpoint, prefix) {
		tlsConfig, err := createTLSConfig(nomadCACert, nomadTLSSkipVerify)
		if err != nil {
			return nil, err
		}

		dpHTTPClient.HTTPClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}
	return  &Nomad{
		Client: &dpHTTPClient,
		URL:    nomadEndpoint,
		// Name:   name,
	}, nil
}

func createTLSConfig (nomadCACert string, nomadTLSSkipVerify bool) (*tls.Config, error) {

	if nomadCACert == "" {
		if !nomadTLSSkipVerify {
			return nil, errors.New("invalid configuration with https but no CA cert or skip verification enabled")
		}
		// no CA file => do not check cert  XXX DANGER DANGER XXX
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

	// assert: nomadCACert is not empty

	// Set caCertPool if first use
	if caCertPool == nil {
		var err error
		caCertPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		if caCertPool == nil {
			caCertPool = x509.NewCertPool()
		}
	}

	caCert, err := ioutil.ReadFile(nomadCACert)
	if err != nil {
		return nil, err
	}
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to append ca cert to pool")
	}

	return &tls.Config{
		RootCAs: caCertPool,
	}, nil

}

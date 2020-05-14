package nomadclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	dphttp "github.com/ONSdigital/dp-net/http"
	"io/ioutil"
	"net/http"
	"strings"
)

type Nomad struct {
	Client dphttp.Clienter
	URL string
	// Name needed for checker later.
	// Name string
}

const prefix = "https://"
var caCertPool *x509.CertPool

func NewClient(nomadEndpoint, nomadCACert string, nomadTLSSkipVerify bool) (*Nomad, error){

	dpHTTPClient := *dphttp.DefaultClient

	if strings.HasPrefix(nomadEndpoint, prefix) {
		tlsConfig, err := createTlsConfig(nomadCACert, nomadTLSSkipVerify)
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

func createTlsConfig (nomadCACert string, nomadLSSkipVerify bool) (*tls.Config, error) {

	if nomadCACert == "" && !nomadLSSkipVerify {
		return nil, errors.New("invalid configuration with https but no CA cert or skip verification enabled")
	}
	if nomadCACert != "" {
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
	// no CA file => do not check cert  XXX DANGER DANGER XXX
	return &tls.Config{
		InsecureSkipVerify: true,
	}, nil
}

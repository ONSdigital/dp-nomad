package nomadclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
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

func NewClient(NomadEndpoint, NomadCACert string, NomadTLSSkipVerify bool) (*Nomad, error){

	dpHTTPClient := *dphttp.DefaultClient

	if strings.HasPrefix(NomadEndpoint, "https://") {
		fmt.Print("in if https")
		var tlsConfig *tls.Config
		if NomadCACert != "" {
			fmt.Print("in nomad CA if")
			caCertPool, _ := x509.SystemCertPool()
			if caCertPool == nil {
				caCertPool = x509.NewCertPool()
			}

			caCert, err := ioutil.ReadFile(NomadCACert)
			if err != nil {
				return nil, err
			}
			if !caCertPool.AppendCertsFromPEM(caCert) {
				return nil, errors.New("failed to append ca cert to pool")
			}

			tlsConfig = &tls.Config{
				RootCAs: caCertPool,
			}
		} else if NomadTLSSkipVerify {
			fmt.Print("in else if TLS")

			// no CA file => do not check cert  XXX DANGER DANGER XXX
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		} else {
			return nil, errors.New("invalid configuration with https but no CA cert or skip verification enabled")
		}
		dpHTTPClient.HTTPClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}
	return  &Nomad{
		Client: &dpHTTPClient,
		URL:    NomadEndpoint,
		// Name:   name,
	}, nil
} 

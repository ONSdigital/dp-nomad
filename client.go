package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/log.go/v2/log"
)

// Nomad represents the nomad client
type Client struct {
	Client dphttp.Clienter
	URL    string
}

// ErrInvalidAppResponse is returned when an app does not respond
// with a valid status
type ErrInvalidAppResponse struct {
	ExpectedCode int
	ActualCode   int
	URI          string
}

const httpsPrefix = "https://"

var caCertPool *x509.CertPool

// Error should be called by the user to print out the stringified version of the error
func (e ErrInvalidAppResponse) Error() string {
	return fmt.Sprintf("invalid response from downstream service - should be: %d, got: %d, path: %s",
		e.ExpectedCode,
		e.ActualCode,
		e.URI,
	)
}

// NewClient returns a Nomad HTTP client for this endpoint
// with optional TLS config
func NewClient(nomadEndpoint, nomadCACert string, nomadTLSSkipVerify bool) (*Client, error) {

	var dpHTTPClient dphttp.Clienter

	if strings.HasPrefix(nomadEndpoint, httpsPrefix) {
		tlsConfig, err := createTLSConfig(nomadCACert, nomadTLSSkipVerify)
		if err != nil {
			return nil, err
		}

		dpHTTPClient = dphttp.NewClientWithTransport(&http.Transport{TLSClientConfig: tlsConfig})
	} else {
		dpHTTPClient = dphttp.NewClient()
	}

	return &Client{
		Client: dpHTTPClient,
		URL:    nomadEndpoint,
	}, nil
}

func createTLSConfig(nomadCACert string, nomadTLSSkipVerify bool) (*tls.Config, error) {

	if nomadCACert == "" {
		if !nomadTLSSkipVerify {
			return nil, errors.New("invalid configuration with https but no CA cert or skip verification enabled")
		}
		// no CA file => do not check cert  XXX DANGER DANGER XXX
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

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

func (c *Client) Get(ctx context.Context, path string) (int, error) {

	req, err := http.NewRequest("GET", c.URL+path, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.Client.Do(ctx, req)
	if err != nil {
		return 0, err
	}
	defer closeResponseBody(ctx, resp)

	if resp.StatusCode < 200 || (resp.StatusCode > 399 && resp.StatusCode != 429) {
		return resp.StatusCode, ErrInvalidAppResponse{http.StatusOK, resp.StatusCode, req.URL.Path}
	}

	return resp.StatusCode, nil
}

func closeResponseBody(ctx context.Context, resp *http.Response) {
	if resp.Body == nil {
		return
	}

	if err := resp.Body.Close(); err != nil {
		log.Error(ctx, "error closing http response body", err)
	}
}

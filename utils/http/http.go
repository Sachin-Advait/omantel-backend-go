package http

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// dbClient is the default client
type httpClient struct {
	*http.Client
}

var client *httpClient

// InitDatabase is used to initialise the database
func InitHttpClient(conf *models.HttpClientConfig) {

	dialer := &net.Dialer{
		Timeout:   conf.ConnectionTimeout,
		KeepAlive: conf.KeepAliveDuration,
	}

	caCert, err := ioutil.ReadFile(conf.CertificatePath)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client = &httpClient{
		&http.Client{
			Transport: &http.Transport{
				DialContext: dialer.DialContext,
				TLSClientConfig: &tls.Config{
					//RootCAs: caCertPool,
					InsecureSkipVerify: true,
				},
				MaxIdleConns:          conf.MaxIdleConnections,
				IdleConnTimeout:       conf.IdleConnectionTimeout,
				TLSHandshakeTimeout:   conf.TLSHandshakeTimeout,
				ExpectContinueTimeout: conf.ExpectContinueTimeout,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			},
			Timeout: conf.Timeout,
		},
	}
}

// Get is used to get the http instance client
func Get() HttpClient {
	return client
}

package httpclient

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

const (
	MaxIdleConnections int = 20
)

type IHTTPClient interface {
	Post(url string, body io.Reader, options Options) (*http.Response, error)
}

type HTTPClient struct {
	httpClient *http.Client
}

type Options struct {
	Headers            map[string]string
	Timeout            time.Duration
	InsecureSkipVerify bool
}

func NewHTTPClient() *HTTPClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
	}

	return &HTTPClient{httpClient}
}

func (h *HTTPClient) setHTTPClientOptions(options Options) {
	h.httpClient.Timeout = options.Timeout
	transport := h.httpClient.Transport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: options.InsecureSkipVerify}
}

func (h *HTTPClient) Post(url string, body io.Reader, options Options) (*http.Response, error) {
	h.setHTTPClientOptions(options)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	return h.httpClient.Do(req)
}

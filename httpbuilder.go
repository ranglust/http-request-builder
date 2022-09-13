package http_request_builder

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpRequestBuilder struct {
	method       string
	username     string
	password     string
	headers      map[string]string
	useBasicAuth bool
	url          string
	body         io.Reader
	skipSSL      bool
}

func NewHTTPRequest() *HttpRequestBuilder {
	return &HttpRequestBuilder{headers: make(map[string]string)}
}

func (h *HttpRequestBuilder) GET() *HttpRequestBuilder {
	h.method = "GET"
	return h
}

func (h *HttpRequestBuilder) DELETE() *HttpRequestBuilder {
	h.method = "DELETE"
	return h
}

func (h *HttpRequestBuilder) PATCH() *HttpRequestBuilder {
	h.method = "PATCH"
	return h
}

func (h *HttpRequestBuilder) POST() *HttpRequestBuilder {
	h.method = "POST"
	return h
}

func (h *HttpRequestBuilder) PUT() *HttpRequestBuilder {
	h.method = "PUT"
	return h
}

func (h *HttpRequestBuilder) WithBasicAuth(username, password string) *HttpRequestBuilder {
	h.username = username
	h.password = password
	h.useBasicAuth = true
	return h
}

func (h *HttpRequestBuilder) AddHeader(key string, value string) *HttpRequestBuilder {
	h.headers[key] = value
	return h
}

func (h *HttpRequestBuilder) SkipSSL() *HttpRequestBuilder {
	h.skipSSL = true
	return h
}

func (h *HttpRequestBuilder) Body(b []byte) *HttpRequestBuilder {
	h.body = bytes.NewBuffer(b)
	return h
}

func (h *HttpRequestBuilder) URL(url string) *HttpRequestBuilder {
	h.url = url
	return h
}

func (h *HttpRequestBuilder) Do() ([]byte, int, error) {
	request, err := http.NewRequest(h.method, h.url, h.body)
	if err != nil {
		return nil, 0, err
	}

	for header, value := range h.headers {
		request.Header.Set(header, value)
	}

	if h.useBasicAuth {
		request.SetBasicAuth(h.username, h.password)
	}

	client := &http.Client{}

	if h.skipSSL {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, resp.StatusCode, nil
}

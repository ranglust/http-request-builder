package http_request_builder

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestHttpRequestBuilder_Negative(t *testing.T) {
	hr := NewHTTPRequest()
	_, _, err := hr.URL("http://localhost:9111").
		GET().
		Do()

	assert.NotNil(t, err)
}
func TestHttpRequestBuilder_Positive(t *testing.T) {
	s, _ := tempHttpServer()
	defer s.Close()

	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()

	hr := NewHTTPRequest()
	body, rc, err := hr.WithBasicAuth("username", "password").
		URL("http://localhost:9111").
		WithParam("param1", "paramvalue1").
		WithHeader("header1", "headervalue1").
		SkipSSL().
		Body([]byte("this is body")).
		POST().
		Do()

	assert.Nil(t, err)
	assert.Equal(t, rc, 200)
	assert.Equal(t, body, []byte(strconv.Itoa(2)))
}

func tempHttpServer() (*http.Server, error) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		count := 0
		header1 := r.Header.Get("header1")
		if header1 == "headervalue1" {
			count++
		}
		params := strings.Split(strings.ReplaceAll(r.RequestURI, "/?", ""), "&")
		for _, p := range params {
			keyValues := strings.Split(p, "=")
			if keyValues[0] == "param1" && keyValues[1] == "paramvalue1" {
				count++
			}
		}
		w.Write([]byte(strconv.Itoa(count)))
	})

	httpServer := http.Server{
		Addr: ":9111",
	}

	return &httpServer, nil
}

package gurl

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestExecute(t *testing.T) {
	tt := []struct {
		name                 string
		tls                  bool
		expectedResponseBody string
		expectedRequestBody  string
		err                  string
		requestHeaders       http.Header
		method               string
		callback             func(c *Config)
	}{
		{
			name:                 "making a request to a server",
			method:               http.MethodGet,
			expectedResponseBody: "Hello, client\n",
			requestHeaders: map[string][]string{
				"Accept-Encoding": {
					"gzip",
				},
				"User-Agent": {
					"gurl",
				},
			},
		},
		{
			name:                 "making a POST request to a server",
			method:               http.MethodPost,
			expectedRequestBody:  "email=user%40example.com&name=user",
			expectedResponseBody: "Hello, client\n",
			requestHeaders: map[string][]string{
				"Accept-Encoding": {
					"gzip",
				},
				"User-Agent": {
					"gurl",
				},
				"Content-Type": {
					"application/x-www-form-urlencoded",
				},
				"Content-Length": {
					"34",
				},
			},
			callback: func(c *Config) {
				params := url.Values{}
				params.Add("name", "user")
				params.Add("email", "user@example.com")
				c.Data = params.Encode()

				c.Method = http.MethodPost

				c.Headers.Add("Content-Type", "application/x-www-form-urlencoded")
			},
		},
		{
			name:                 "making a request to a tls server with insecure on",
			method:               http.MethodGet,
			expectedResponseBody: "Hello, client\n",
			requestHeaders: map[string][]string{
				"Accept-Encoding": {
					"gzip",
				},
				"User-Agent": {
					"gurl",
				},
			},
			tls: true,
			callback: func(c *Config) {
				c.Insecure = true
			},
		},
		{
			name:                 "making a request to a tls server with insecure off",
			method:               http.MethodGet,
			expectedResponseBody: "Hello, client\n",
			tls:                  true,
			err:                  "x509: “Acme Co” certificate is not trusted",
			requestHeaders: map[string][]string{
				"Accept-Encoding": {
					"gzip",
				},
				"User-Agent": {
					"gurl",
				},
			},
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			var server *httptest.Server

			var receivedHeaders http.Header
			var requestBody string
			var method string

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedHeaders = r.Header
				method = r.Method

				if bodyBytes, err := ioutil.ReadAll(r.Body); err == nil {
					requestBody = string(bodyBytes)
				}

				fmt.Fprintln(w, "Hello, client")
			})

			if ts.tls {
				server = httptest.NewTLSServer(handler)
			} else {
				server = httptest.NewServer(handler)
			}

			defer server.Close()

			responseBody := bytes.NewBuffer(make([]byte, 0, 1024))

			c := &Config{
				Headers:            map[string][]string{},
				UserAgent:          "gurl",
				Method:             http.MethodGet,
				Url:                mustURL(t, server.URL),
				ControlOutput:      bytes.NewBuffer(make([]byte, 0, 1024)),
				ResponseBodyOutput: responseBody,
			}

			if ts.callback != nil {
				ts.callback(c)
			}

			err := Execute(c)
			if ts.err != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), ts.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, ts.requestHeaders, receivedHeaders)
				assert.Equal(t, ts.expectedRequestBody, requestBody)
				assert.Equal(t, ts.expectedResponseBody, responseBody.String())
				assert.Equal(t, ts.method, method)
			}
		})
	}
}

package gurl

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Headers            http.Header
	UserAgent          string
	Data               string
	Method             string
	Insecure           bool
	Url                *url.URL
	ControlOutput      io.Writer
	ResponseBodyOutput io.Writer
}

func Execute(c *Config) error {
	var r io.Reader
	var tlsConfig *tls.Config

	if c.Data != "" {
		r = bytes.NewBufferString(c.Data)
	}

	if c.Insecure {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	request, err := http.NewRequest(c.Method, c.Url.String(), r)
	if err != nil {
		return err
	}

	if c.UserAgent != "" {
		request.Header.Set("User-Agent", c.UserAgent)
	}

	for key, values := range c.Headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	requestBuilder := &wrappedBuilder{
		prefix: ">",
	}

	requestBuilder.Printf("%v %v", request.Method, request.URL.String())
	requestBuilder.WriteHeaders(request.Header)
	requestBuilder.Println()

	if _, err := io.Copy(c.ControlOutput, strings.NewReader(requestBuilder.String())); err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Warn().Err(err).Str("url", c.Url.String()).Msg("failed to close response body")
		}
	}()

	responseBuilder := &wrappedBuilder{
		prefix: "<",
	}

	responseBuilder.Printf("%v %v", response.Proto, response.Status)
	responseBuilder.WriteHeaders(response.Header)
	responseBuilder.Printf("")
	responseBuilder.Println()

	if _, err := io.Copy(c.ControlOutput, strings.NewReader(responseBuilder.String())); err != nil {
		return err
	}

	_, err = io.Copy(c.ResponseBodyOutput, response.Body)
	return err
}

type wrappedBuilder struct {
	prefix string
	strings.Builder
}

func (w *wrappedBuilder) WriteHeaders(headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			w.Printf("%v: %v", key, value)
		}
	}
}

func (w *wrappedBuilder) Println() {
	w.WriteString("\n")
}

func (w *wrappedBuilder) Printf(s string, a ...any) {
	w.WriteString(fmt.Sprintf("%v %v\n", w.prefix, fmt.Sprintf(s, a...)))
}

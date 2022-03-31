package gurl

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func mustURL(t *testing.T, u string) *url.URL {
	result, err := url.Parse(u)
	require.NoError(t, err)

	return result
}

func TestArgsValidator(t *testing.T) {
	tt := []struct {
		name           string
		args           []string
		err            string
		expectedConfig *Config
	}{
		{
			name: "a good execution with a valid url",
			args: []string{
				"https://example.com/request",
			},
			expectedConfig: &Config{
				Url: mustURL(t, "https://example.com/request"),
			},
		},
		{
			name: "an invalid url",
			args: []string{
				"not a valid url\t",
			},
			err: "the URL provided is invalid: not a valid url\t: parse \"not a valid url\\t\": net/url: invalid control character in URL",
		},
		{
			name: "no url",
			args: []string{},
			err:  "you must provide a single URL to be called but you provided 0",
		},
		{
			name: "too many urls",
			args: []string{
				"https://example,com/request",
				"https://example.com/response",
			},
			err: "you must provide a single URL to be called but you provided 2",
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			c := &Config{}

			err := ArgsValidator(c)(nil, ts.args)
			if ts.err != "" {
				assert.EqualError(t, err, ts.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, ts.expectedConfig, c)
			}
		})
	}
}

func TestOptionsValidator(t *testing.T) {
	tt := []struct {
		name           string
		expectedConfig *Config
		headers        []string
		err            string
	}{
		{
			name: "with valid headers",
			headers: []string{
				"User-Agent: Sample",
				"Cache: None",
				"Cache: Respect",
			},
			expectedConfig: &Config{
				Headers: map[string][]string{
					"User-Agent": {
						"Sample",
					},
					"Cache": {
						"None",
						"Respect",
					},
				},
			},
		},
		{
			name: "with invalid headers",
			headers: []string{
				"User-Agent",
			},
			err: "header is not a valid http header separated by `:`, value was: [User-Agent]",
		},
	}

	for _, ts := range tt {
		t.Run(ts.name, func(t *testing.T) {
			c := &Config{
				Headers: map[string][]string{},
			}
			err := OptionsValidator(c, ts.headers)(nil, nil)
			if ts.err != "" {
				assert.EqualError(t, err, ts.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, ts.expectedConfig, c)
			}
		})
	}

}

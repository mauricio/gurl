package gurl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func CreateCommand() *cobra.Command {
	config := &Config{
		Headers:            map[string][]string{},
		ResponseBodyOutput: os.Stdout,
		ControlOutput:      os.Stdout,
	}

	headers := make([]string, 0, 255)

	command := &cobra.Command{
		Use:   `gurl [options] URL`,
		Short: `gurl is an HTTP client`,
		Long:  `gurl is an HTTP client for a tutorial on how to build command line clients in go`,
		Args: func(cmd *cobra.Command, args []string) error {
			if l := len(args); l != 1 {
				return fmt.Errorf("you must provide a single URL to be called but you provided %v", l)
			}

			u, err := url.Parse(args[0])
			if err != nil {
				return errors.Wrapf(err, "the URL provided is invalid: %v", args[0])
			}

			config.Url = u

			return nil
		},
		PreRunE: PreRun(config, headers),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Execute(config)
		},
	}

	command.LocalFlags().StringSliceVarP(&headers, "header", "h", nil, `custom header to be sent to the server, can be set multiple times, headers are set as "HeaderName: Header content"`)
	command.LocalFlags().StringVarP(&config.UserAgent, "user-agent", "u", "gurl", "the user agent to be used for requests")
	command.LocalFlags().StringVarP(&config.Data, "data", "d", "", "data to be sent as the request body")
	command.LocalFlags().StringVarP(&config.Method, "method", "m", http.MethodGet, "HTTP method to be used for the request, defaults to GET")
	command.LocalFlags().BoolVarP(&config.Insecure, "insecure", "k", false, "allows insecure server connections over HTTPS")

	return command
}

func PreRun(c *Config, headers []string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, h := range headers {
			if name, value, found := strings.Cut(h, ":"); found {
				c.Headers.Add(name, value)
			} else {
				return errors.Errorf("header is not a valid http header separated by `:`, value was: [%v]", h)
			}
		}

		return nil
	}
}

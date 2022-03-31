package gurl

import (
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
		Use:     `gurl URL`,
		Short:   `gurl is an HTTP client`,
		Long:    `gurl is an HTTP client for a tutorial on how to build command line clients in go`,
		Args:    ArgsValidator(config),
		PreRunE: OptionsValidator(config, headers),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Execute(config)
		},
	}

	command.PersistentFlags().StringSliceVarP(&headers, "headers", "H", nil, `custom headers headers to be sent with the request, headers are separated by "," as in "HeaderName: Header content,OtherHeader: Some other value"`)
	command.PersistentFlags().StringVarP(&config.UserAgent, "user-agent", "u", "gurl", "the user agent to be used for requests")
	command.PersistentFlags().StringVarP(&config.Data, "data", "d", "", "data to be sent as the request body")
	command.PersistentFlags().StringVarP(&config.Method, "method", "m", http.MethodGet, "HTTP method to be used for the request")
	command.PersistentFlags().BoolVarP(&config.Insecure, "insecure", "k", false, "allows insecure server connections over HTTPS")

	return command
}

func ArgsValidator(c *Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if l := len(args); l != 1 {
			return newErrorWithCode(2, "you must provide a single URL to be called but you provided %v", l)
		}

		u, err := url.Parse(args[0])
		if err != nil {
			return errors.Wrapf(err, "the URL provided is invalid: %v", args[0])
		}

		c.Url = u

		return nil
	}
}

func OptionsValidator(c *Config, headers []string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, h := range headers {
			if name, value, found := strings.Cut(h, ":"); found {
				c.Headers.Add(strings.TrimSpace(name), strings.TrimSpace(value))
			} else {
				return newErrorWithCode(3, "header is not a valid http header separated by `:`, value was: [%v]", h)
			}
		}

		return nil
	}
}

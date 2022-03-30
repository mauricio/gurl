package gurl

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/url"
)

type Config struct {
	Headers   map[string][]string
	UserAgent string
}

func CreateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   `gurl [options] URL`,
		Short: `gurl is an HTTP client`,
		Long:  `gurl is an HTTP client for a tutorial on how to build command line clients in go`,
		Args: func(cmd *cobra.Command, args []string) error {
			if l := len(args); l != 1 {
				return errors.New("you must provide a single URL to be called")
			}

			_, err := url.Parse(args[0])
			if err != nil {
				return errors.Wrapf(err, "the URL provided is invalid: %v", args[0])
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {

		},
	}

	headers := make([]string, 0, 10)

	command.LocalFlags().StringArrayVarP(&headers, "header", "h", nil, `custom header to be sent to the server, can be set multiple times, headers are set as "HeaderName: Header content"`)

	return command
}

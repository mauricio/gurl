package main

import (
	"github.com/mauricio/gurl"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := gurl.CreateCommand().Execute(); err != nil {
		switch e := err.(type) {
		case gurl.ReturnCodeError:
			os.Exit(e.Code())
		default:
			os.Exit(1)
		}
	}
}

package main

import (
	"github.com/mauricio/gurl"
	"os"
)

func main() {
	if err := gurl.CreateCommand().Execute(); err != nil {
		switch e := err.(type) {
		case gurl.ReturnCodeError:
			os.Exit(e.Code())
		default:
			os.Exit(1)
		}
	}
}

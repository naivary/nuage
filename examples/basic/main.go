package main

import (
	"fmt"
	"os"

	"github.com/naivary/nuage"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stdout, "err: %v", err)
	}
}

func run() error {
	api, err := nuage.New()
	if err != nil {
		return err
	}
	_ = api
	return nil
}

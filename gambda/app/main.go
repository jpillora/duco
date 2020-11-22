package main

import (
	"gambda/runtime"
	"io"
	"log"
)

func main() {
	r := runtime.New()

	r.HandleFunc("myhandler", func(input io.Reader, output io.Writer) error {
		_, err := io.Copy(output, input)
		return err
	})

	if err := r.Start(); err != nil {
		log.Fatal(err)
	}
}

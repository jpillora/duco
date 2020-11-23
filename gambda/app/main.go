package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gambda/runtime"
	"io"
	"io/ioutil"
	"log"
)

func main() {
	r := runtime.New()

	r.HandleFunc("myhandler", func(ctx runtime.Context, input io.Reader, output io.Writer) error {
		b, err := ioutil.ReadAll(input)
		if err != nil {
			return err
		}
		if bytes.Contains(b, []byte("err")) {
			return fmt.Errorf("it errored '%s'", b)
		}

		var client interface{}
		ctx.DecodeClient(&client)

		e := json.NewEncoder(output)
		e.SetIndent("", "  ")
		return e.Encode(map[string]interface{}{
			"traceid":   ctx.TraceID(),
			"cogid":     ctx.CognitoIdentity(),
			"deadline":  ctx.FunctionDeadline(),
			"clientctx": client,
			"input":     string(b),
		})
	})

	if err := r.Start(); err != nil {
		log.Fatal(err)
	}
}

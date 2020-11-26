package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gambda"
	"io"
	"io/ioutil"
)

func New() gambda.Func {
	return &myhandler{}
}

type myhandler struct {
}

func (h *myhandler) Name() string {
	return "myhandler"
}

func (h *myhandler) Handle(ctx gambda.Context, input io.Reader, output io.Writer) error {
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
		"arn":       ctx.InvokedFunctionARN(),
	})
}

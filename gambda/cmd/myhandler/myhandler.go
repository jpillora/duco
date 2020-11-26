package myhandler

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

	var v interface{}
	ctx.DecodeClient(&v)

	e := json.NewEncoder(output)
	e.SetIndent("", "  ")
	return e.Encode(map[string]interface{}{
		"trace-id": ctx.TraceID(),
		"cog-id":   ctx.CognitoIdentity(),
		"deadline": ctx.FunctionDeadline(),
		"client":   v,
		"input":    string(b),
		"arn":      ctx.InvokedFunctionARN(),
	})
}

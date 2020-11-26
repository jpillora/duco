package duco

import (
	"context"
	"io"
	"time"
)

//Context is a duco context
type Context interface {
	//Runtime context
	context.Context
	//Lambda context
	FunctionDeadline() time.Time
	DecodeClient(v interface{}) error
	TraceID() string
	CognitoIdentity() string
	InvokedFunctionARN() string
}

type NewFunc func() Func

type Func interface {
	Name() string
	Handle(Context, io.Reader, io.Writer) error
}

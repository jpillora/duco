package gambda

import (
	"context"
	"io"
	"time"
)

//Context is a gambda context
type Context interface {
	//Invoke context
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

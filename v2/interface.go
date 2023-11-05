package duco

import (
	"context"
	"io"
	"time"
)

//Context is a duco context,
//which is active for the duration of the invocation,
//and also provides "lambda context" fields.
type Context interface {
	context.Context
	FunctionDeadline() time.Time
	DecodeClient(v interface{}) error
	TraceID() string
	CognitoIdentity() string
	InvokedFunctionARN() string
}

//Func is a lambda function handler
type Func interface {
	Handle(Context, io.Reader, io.Writer) error
}

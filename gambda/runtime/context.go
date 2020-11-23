package runtime

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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

type gcontext struct {
	inner   context.Context
	headers http.Header
}

func (c *gcontext) FunctionDeadline() time.Time {
	if n, err := strconv.ParseInt(c.headers.Get("Lambda-Runtime-Deadline-Ms"), 10, 64); err == nil {
		return time.Unix(n/1e3, (n%1e3)*1e6)
	}
	return time.Time{}
}

func (c *gcontext) TraceID() string {
	return c.headers.Get("Lambda-Runtime-Trace-Id")
}

func (c *gcontext) InvokedFunctionARN() string {
	return c.headers.Get("Lambda-Runtime-Invoked-Function-Arn")
}

func (c *gcontext) CognitoIdentity() string {
	return c.headers.Get("Lambda-Runtime-Cognito-Identity")
}

func (c *gcontext) DecodeClient(v interface{}) error {
	return json.Unmarshal(
		[]byte(c.headers.Get("Lambda-Runtime-Client-Context")),
		v,
	)
}

//proxy context.Context
func (c *gcontext) Deadline() (deadline time.Time, ok bool) {
	return c.inner.Deadline()
}

func (c *gcontext) Done() <-chan struct{} {
	return c.inner.Done()
}

func (c *gcontext) Err() error {
	return c.inner.Err()
}

func (c *gcontext) Value(key interface{}) interface{} {
	return c.inner.Value(key)
}

package duco

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type invocation struct {
	//params:
	r       *Runtime
	n       int
	headers http.Header
	fnIn    io.ReadCloser
	//computed:
	fnOut     io.WriteCloser
	respIn    io.ReadCloser
	written   bool
	responded bool
}

func (i *invocation) id() string {
	return i.headers.Get("Lambda-Runtime-Aws-Request-Id")
}

func (i *invocation) logf(format string, args ...interface{}) {
	if !i.r.debug {
		return
	}
	id := i.id()
	prefix := fmt.Sprintf("[invocation %s#%d] ", id[0:6], i.n)
	log.Printf(prefix+format, args...)
}

func (i *invocation) handle() {
	//time
	t0 := time.Now()
	i.logf("start")
	defer func() { i.logf("took %s", time.Since(t0)) }()
	//handle only executes while fn blocks,
	//defer cancel communicates this downstream into fn
	inner, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := gcontext{
		inner:   inner,
		headers: i.headers,
	}
	//pipe output of user function,
	//into the bootstrap response
	i.respIn, i.fnOut = io.Pipe()
	//execute function
	fnOutWrapper := io.WriteCloser(i)
	//fn can either
	//error: and respond with a static string
	//success: and respond with a stream
	if err := i.r.fn.Handle(&ctx, i.fnIn, fnOutWrapper); err != nil {
		i.logf("errored: %s", err)
		i.respond(true, strings.NewReader(err.Error()))
	}
	//close pipe which closes response body
	fnOutWrapper.Close()
}

func (i *invocation) Write(b []byte) (int, error) {
	//pipe fnOut to respIn on *first* write
	if !i.written {
		i.written = true
		go i.respond(false, i.respIn)
	}
	return i.fnOut.Write(b)
}

func (i *invocation) Close() error {
	if !i.responded {
		i.respond(false, nil)
	}
	return i.fnOut.Close()
}

func (i *invocation) respond(errd bool, body io.Reader) {
	if i.responded {
		return
	}
	i.responded = true
	action := "response"
	if errd {
		action = "error"
	}
	url := fmt.Sprintf(
		"http://%s/2018-06-01/runtime/invocation/%s/%s",
		i.r.env.api, i.id(), action,
	)
	resp, err := http.Post(url, "application/octet-stream", body)
	if err != nil {
		panic(err)
	} else {
		i.logf("%s %s", action, http.StatusText(resp.StatusCode))
	}
}

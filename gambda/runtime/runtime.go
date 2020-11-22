package runtime

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func New() *Runtime {
	return &Runtime{
		env: loadEnv(),
		fns: map[string]Func{},
	}
}

type Runtime struct {
	env
	fns map[string]Func
}

func (r *Runtime) HandleFunc(name string, fn Func) {
	r.fns[name] = fn
}

func (r *Runtime) Start() error {
	if r.env.dev {
		return r.startDevelopment()
	}
	return r.startLambda()
}

func (r *Runtime) startDevelopment() error {
	panic("no implemented")
}

func (r *Runtime) startLambda() error {
	for {
		log.Printf("invoke next...")
		if err := r.invokeNext(); err != nil {
			return r.invokeNext()
		}
		time.Sleep(time.Second)
	}
}

func (r *Runtime) invokeNext() error {
	resp, err := http.Get(fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/next", r.env.api))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	reqID := resp.Header.Get("Lambda-Runtime-Aws-Request-Id")
	log.Printf("[invocation %s] start", reqID)

	fn, ok := r.fns["myhandler"]
	if !ok {
		panic("no func")
	}

	eg := errgroup.Group{}

	pr, pw := io.Pipe()

	eg.Go(func() error {
		input := resp.Body
		output := pw
		err := fn(input, output)
		pw.Close()
		log.Printf("[invocation %s] handled %v", reqID, err)
		return err
	})

	eg.Go(func() error {
		url := fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/%s/response", r.env.api, reqID)
		resp, err := http.Post(url, "application/octet-stream", pr)
		if err != nil {
			return err
		}
		log.Printf("[invocation %s] responded -> %s", reqID, http.StatusText(resp.StatusCode))
		return nil
	})

	return eg.Wait()
}

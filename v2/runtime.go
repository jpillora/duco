package duco

import (
	"fmt"
	"net/http"
	"strings"
)

func New(fn Func) *Runtime {
	r := &Runtime{
		env: loadEnv(),
		fn:  fn,
	}
	return r
}

type Runtime struct {
	env
	fn      Func
	invokes int
	debug   bool
}

//Debug enables logs
func (r *Runtime) Debug() *Runtime {
	r.debug = true
	return r
}

//Start the duco runtime
func (r *Runtime) Start() error {
	err := r.start()
	if err != nil {
		//tell lambda hypervisor about this...
		url := fmt.Sprintf("http://%s/2018-06-01/runtime/init/error", r.env.api)
		http.Post(url, "application/octet-stream", strings.NewReader(err.Error()))
	}
	return err
}

func (r *Runtime) start() error {
	if r.env.dev {
		return r.startDevelopment()
	}
	return r.startLambda()
}

func (r *Runtime) startDevelopment() error {
	// log.Printf("listening on 8081...")
	// return http.ListenAndServe(":8081", wrapFunc(r.fn))
	panic("TODO")
}

func (r *Runtime) startLambda() error {
	for {
		if err := r.invokeNext(); err != nil {
			return err
		}
		r.invokes++
	}
}

func (r *Runtime) invokeNext() error {
	//blocks here while we wait for next request
	resp, err := http.Get(fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/next", r.env.api))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//invoke!
	inv := invocation{
		r:       r,
		n:       r.invokes,
		headers: resp.Header,
		fnIn:    resp.Body,
	}
	if inv.id() == "" {
		return nil //skip no id...
	}
	inv.handle()
	return nil
}

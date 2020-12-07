package runtime

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jpillora/duco"
)

func New() *Runtime {
	return &Runtime{
		env: loadEnv(),
		fns: map[string]duco.Func{},
	}
}

type Runtime struct {
	env
	fns     map[string]duco.Func
	invokes int
}

func (r *Runtime) Add(fn duco.Func) {
	r.fns[fn.Name()] = fn
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
	log.Printf("listening on 8081...")
	return http.ListenAndServe(":8081", nil)
}

func (r *Runtime) fn() duco.Func {
	for _, fn := range r.fns {
		return fn
	}
	fn, ok := r.fns["myhandler"]
	if !ok {
		panic("no func")
	}
	return fn
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
	fn := i.r.fn()
	fnOutWrapper := io.WriteCloser(i)
	//fn can either
	//error: and respond with a static string
	//success: and respond with a stream
	if err := fn.Handle(&ctx, i.fnIn, fnOutWrapper); err != nil {
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

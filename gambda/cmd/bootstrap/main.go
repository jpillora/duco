package main

import (
	"gambda"
	"gambda/runtime"
	"io/ioutil"
	"log"
	"plugin"
)

func main() {
	r := runtime.New()

	ss, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range ss {
		p, err := plugin.Open(s.Name())
		if err != nil {
			log.Fatalf("plugin-open: %s", err)
		}
		log.Printf("openned %s", s.Name())
		v, err := p.Lookup("New")
		if err != nil {
			log.Fatalf("plugin-lookup: %s", err)
		}
		newFn, ok := v.(func() gambda.Func)
		if !ok {
			log.Fatalf("new-fn: got %T", newFn)
		}
		f := (newFn)()
		r.HandleFunc(f)
	}

	if err := r.Start(); err != nil {
		log.Fatal(err)
	}

}

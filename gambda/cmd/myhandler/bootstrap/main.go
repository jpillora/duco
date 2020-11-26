package main

import (
	"gambda/cmd/myhandler"
	"gambda/runtime"
)

func main() {
	g := runtime.New()
	g.Add(myhandler.New())
	g.Start()
}

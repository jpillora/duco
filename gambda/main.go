package main

import (
	"gambda/command/deploy"
	"gambda/command/invoke"

	"github.com/jpillora/opts"
)

func main() {
	opts.
		New(&struct{}{}).
		Name("gambda").
		AddCommand(deploy.Command()).
		AddCommand(invoke.Command()).
		Parse().
		RunFatal()
}

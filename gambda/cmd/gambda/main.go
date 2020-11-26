package main

import (
	"gambda/cmd/gambda/deploy"
	"gambda/cmd/gambda/invoke"
	"gambda/cmd/gambda/layer"

	"github.com/jpillora/opts"
)

func main() {
	opts.
		New(&struct{}{}).
		Name("gambda").
		AddCommand(deploy.Command()).
		AddCommand(invoke.Command()).
		AddCommand(layer.Command()).
		Complete().
		Parse().
		RunFatal()
}

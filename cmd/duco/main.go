package main

import (
	"github.com/jpillora/duco/internal/cmd/deploy"
	"github.com/jpillora/duco/internal/cmd/invoke"
	"github.com/jpillora/duco/internal/cmd/layer"

	"github.com/jpillora/opts"
)

func main() {
	opts.
		New(&struct{}{}).
		Name("duco").
		Repo("github.com/jpillora/duco").
		AddCommand(deploy.Command()).
		AddCommand(invoke.Command()).
		AddCommand(layer.Command()).
		Complete().
		Parse().
		RunFatal()
}

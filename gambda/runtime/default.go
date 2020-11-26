package runtime

import "gambda"

var DefaultRuntime = New()

func HandleFunc(fn gambda.Func) {
	DefaultRuntime.HandleFunc(fn)
}

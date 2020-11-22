package runtime

import "io"

// type Func interface{}

type Func func(io.Reader, io.Writer) error

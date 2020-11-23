package runtime

import "io"

// type Func interface{}

type Func func(Context, io.Reader, io.Writer) error

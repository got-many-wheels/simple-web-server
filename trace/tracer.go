package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface that describees an object capable of
// tracing events throughout code
type Tracer interface {
	Trace(...any)
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...any) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// a nil tracer that suppoesed to do nothing
type nilTracer struct{}

func (t *nilTracer) Trace(a ...any) {}

func Off() Tracer {
	return &nilTracer{}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

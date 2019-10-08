package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	return withStackN(err, 4)
}

func withStackInner(err error) error {
	return withStackN(err, 5)
}

func withStackN(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(skip),
	}
}

type withStack struct {
	err   error
	stack *stack
}

func (w *withStack) Error() string { return w.err.Error() }
func (w *withStack) Unwrap() error { return w.err }

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) format() []string {
	var res []string
	frames := runtime.CallersFrames([]uintptr(*s))
	for {
		frame, more := frames.Next()
		res = append(res, fmt.Sprintf("%s:%d   %s", frame.File, frame.Line, funcname(frame.Function)))
		if !more {
			break
		}
	}
	return res
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

func simplifyStack(a []string, b []string) ([]string, int) {
	i := 0
	for i < len(a)-1 && i < len(b)-1 && a[len(a)-i-1] == b[len(b)-i-1] {
		i++
	}

	return a[:len(a)-i], i
}

func callers(skip int) *stack {
	const depth = 96
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func StackTrace(err error) string {
	var b bytes.Buffer

	b.Write([]byte(err.Error()))
	doStackTrace(&b, err, nil)

	return b.String()
}

func doStackTrace(w *bytes.Buffer, err error, next []string) {
	var simplified []string
	var n int

	if err, ok := err.(*withStack); ok {
		raw := err.stack.format()
		simplified, n = simplifyStack(raw, next)
		next = raw
	}

	err2 := Unwrap(err)
	if err != nil {
		doStackTrace(w, err2, next)
	}

	if err, ok := err.(*withStack); ok {
		w.WriteByte('\n')
		w.Write([]byte(err.Error()))
		for _, s := range simplified {
			w.WriteString("\n    ")
			w.WriteString(s)
		}
		if n != 0 {
			fmt.Fprintf(w, "\n    %d more...", n)
		}
	}
}

// multilog is a simple wrapper around github.com/charmbracelet/log that enables creating loggers with multiple outputs.
//
// for example:
//
//	logFile, err := os.Create("log.txt")
//
//	log := multilog.NewMulti(os.Stdout, logFile)
//	log.Info("logging info into stdout and a file!")
//
// The one big issue stemming from this package using the io.Writer interface
// is that multiple errors can not be handled as io.Writer returns a single error.
//
// This means only the first error is returned when charmbracelet/log tries to write the log.
// Making it possibly hard to debug as charmbracelet/log does not provide very good error info out of the io.Writer.
//
// I did consider returning an error channel with the log.Logger
// but that would require the user to deal with handling wirter errors on their own
// and asynchronously, which isn't ideal as I want to keep this as simple as the original charmbracelet/log.
package multilog

import (
	"io"

	"github.com/charmbracelet/log"

	"github.com/kociumba/multilog/multiwriter"
)

// front facing user api

// NewMulti returns a new logger using the default options.
func NewMulti(writers ...io.Writer) *log.Logger {
	w := newMultiWriter(writers...)

	return log.New(w)
}

// NewMultiWithOptions returns a new logger using the provided options.
//
// The writers need to be passed in as the last parameter because of go's limitations around variadic arguments.
func NewMultiWithOptions(o log.Options, writers ...io.Writer) *log.Logger {
	w := newMultiWriter(writers...)

	return log.NewWithOptions(w, o)
}

// package implementation

// HELPER: just initializes and allocates a new MultiWriter
func newMultiWriter(w ...io.Writer) io.Writer {
	return &multiwriter.MultiWriter{Writers: w}
}

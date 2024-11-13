// A small implementation of the io.Writer interface that allows simultaneous writing to multiple io.Writers
//
// Only the first write error is pushed out and all of the writes are executed in parallel.
//
// Usage:
//
//	multi := NewMultiWriter(os.Stdout, buf, logFile, TCPWriter)
//
//	multi.Write([]byte("writing into 4 writers simultaneously!"))
//
// Or you can provide an ErrorHandler callback that will be called for each error encountered during Write:
//
//	multi := NewMultiWriter(os.Stdout, buf, logFile, TCPWriter, failingWriter)
//	multi.ErrorHandler = func(we WriteError) {
//		log.Printf("Write failed for writer %v: %v", we.Writer, we.Err)
//	}
//
//	multi.Write([]byte("writing into 4 writers simultaneously with error handling!"))
package multiwriter

import (
	"io"
	"sync"
)

// MultiWriter allows simultaneous writing to multiple io.Writer interfaces.
type MultiWriter struct {
	writers []io.Writer // the array of writers, all of them will write when the Write method is called
	// ErrorHandler is called for each error encountered during Write.
	// If nil, only the first error is returned from Write.
	ErrorHandler func(WriteError)
}

// WriteError contains a single write error with its associated writer
type WriteError struct {
	Writer io.Writer // the writer that encountered the error
	Err    error     // the error returned by the writer
}

// Write writes the contents of p to each of the underlying io.Writers.
//
// A Write error does not stop any other writes but will be returned !
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	var wg sync.WaitGroup
	wg.Add(len(mw.writers))

	for _, w := range mw.writers {
		go func(w io.Writer) {
			defer wg.Done()

			n, writeErr := w.Write(p)
			if writeErr != nil || n != len(p) {
				if writeErr == nil {
					writeErr = io.ErrShortWrite
				}
				// callback execute
				if mw.ErrorHandler != nil {
					mw.ErrorHandler(WriteError{Writer: w, Err: writeErr})
				}
				if err == nil {
					err = writeErr // Store first error to return
				}
			}
		}(w)
	}

	wg.Wait()
	return len(p), err
}

// initializes and allocates a new MultiWriter
//
// usefull if you want to use the multiwriter by itself without the charm/log wrapper
func NewMultiWriter(w ...io.Writer) io.Writer {
	return &MultiWriter{writers: w}
}

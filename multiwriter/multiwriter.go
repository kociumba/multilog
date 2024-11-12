package multiwriter

import (
	"io"
	"sync"
)

// MultiWriter allows simultaneous writing to multiple io.Writer interfaces.
type MultiWriter struct {
	Writers []io.Writer
}

/*
Write writes the contents of p to each of the underlying io.Writers.

A Write error does not stop any other writes but will be returned !
*/
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	var writeGroup sync.WaitGroup
	errChan := make(chan error)

	for _, writer := range mw.Writers {
		writeGroup.Add(1)

		go func(writer io.Writer) {
			defer writeGroup.Done()
			_, err = writer.Write(p)
			if err != nil {
				errChan <- err
			}
		}(writer)
	}

	go func() {
		writeGroup.Wait()
		close(errChan)
	}()

	select {
	case err = <-errChan:
		return len(p), err
	default:
		return len(p), nil
	}
}

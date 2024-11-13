package multiwriter_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/kociumba/multilog/multiwriter"
)

type failingWriter struct {
	err error
}

func (fw *failingWriter) Write(p []byte) (n int, err error) {
	return 0, fw.err
}

func TestMultiWriter(t *testing.T) {
	t.Run("writes to multiple writers", func(t *testing.T) {
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}
		mw := multiwriter.NewMultiWriter(buf1, buf2)

		testData := []byte("test message")
		n, err := mw.Write(testData)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if n != len(testData) {
			t.Errorf("expected %d bytes written, got %d", len(testData), n)
		}
		if buf1.String() != string(testData) {
			t.Errorf("buffer 1: expected %q, got %q", testData, buf1.String())
		}
		if buf2.String() != string(testData) {
			t.Errorf("buffer 2: expected %q, got %q", testData, buf2.String())
		}
	})

	t.Run("handles writer errors", func(t *testing.T) {
		buf := &bytes.Buffer{}
		expectedErr := errors.New("write error")
		fw := &failingWriter{err: expectedErr}
		mw := multiwriter.NewMultiWriter(buf, fw)

		testData := []byte("test message")
		_, err := mw.Write(testData)

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		// The successful writer should still have written the data
		if buf.String() != string(testData) {
			t.Errorf("buffer: expected %q, got %q", testData, buf.String())
		}
	})

	t.Run("handles error callback", func(t *testing.T) {
		buf := &bytes.Buffer{}
		expectedErr := errors.New("write error")
		fw := &failingWriter{err: expectedErr}
		
		var gotWriter io.Writer
		var gotErr error
		
		mw := multiwriter.NewMultiWriter(buf, fw).(*multiwriter.MultiWriter)
		mw.ErrorHandler = func(we multiwriter.WriteError) {
			gotWriter = we.Writer
			gotErr = we.Err
		}

		testData := []byte("test message")
		_, _ = mw.Write(testData)

		if gotWriter != fw {
			t.Error("error handler received wrong writer")
		}
		if gotErr != expectedErr {
			t.Errorf("error handler: expected error %v, got %v", expectedErr, gotErr)
		}
	})

	t.Run("handles short writes", func(t *testing.T) {
		shortWriter := &shortWriter{maxBytes: 4}
		mw := multiwriter.NewMultiWriter(shortWriter)

		testData := []byte("test message")
		_, err := mw.Write(testData)

		if err != io.ErrShortWrite {
			t.Errorf("expected ErrShortWrite, got %v", err)
		}
	})
}

type shortWriter struct {
	maxBytes int
}

func (sw *shortWriter) Write(p []byte) (n int, err error) {
	if len(p) > sw.maxBytes {
		return sw.maxBytes, nil
	}
	return len(p), nil
}

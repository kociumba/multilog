package multilog_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/kociumba/multilog"
)

func TestNewMulti(t *testing.T) {
	// Create buffers to capture output
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	// Create a new multi logger with explicit options to ensure consistent behavior
	logger := multilog.NewMultiWithOptions(log.Options{
		Level: log.InfoLevel,
		ReportCaller: false,
		ReportTimestamp: false,
		Prefix: "TEST:",
	}, buf1, buf2)

	// Test logging
	testMessage := "test message"
	logger.Info(testMessage)

	// Debug output contents
	t.Logf("Buffer 1 contents: %q", buf1.String())
	t.Logf("Buffer 2 contents: %q", buf2.String())

	// Check if both buffers received the message
	if !bytes.Contains(buf1.Bytes(), []byte(testMessage)) {
		t.Errorf("Buffer 1 did not contain message.\nExpected to contain: %q\nActual contents: %q", testMessage, buf1.String())
	}
	if !bytes.Contains(buf2.Bytes(), []byte(testMessage)) {
		t.Errorf("Buffer 2 did not contain message.\nExpected to contain: %q\nActual contents: %q", testMessage, buf2.String())
	}
}

func TestNewMultiWithOptions(t *testing.T) {
	// Create buffers to capture output
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	// Create custom options
	opts := log.Options{
		Level:           log.DebugLevel,
		ReportCaller:    false,
		ReportTimestamp: false,
		Prefix:          "TEST:",
	}

	// Create a new multi logger with options
	logger := multilog.NewMultiWithOptions(opts, buf1, buf2)

	// Test logging at different levels
	testMessage := "test message"

	// Debug message should appear (level is set to Debug)
	logger.Debug(testMessage)
	t.Logf("After Debug - Buffer 1: %q", buf1.String())
	
	if !bytes.Contains(buf1.Bytes(), []byte(testMessage)) {
		t.Errorf("Debug message not logged to buffer 1.\nBuffer contents: %q", buf1.String())
	}

	// Clear buffers
	buf1.Reset()
	buf2.Reset()

	// Info message should appear
	logger.Info(testMessage)
	t.Logf("After Info - Buffer 1: %q", buf1.String())
	
	if !bytes.Contains(buf1.Bytes(), []byte(testMessage)) {
		t.Errorf("Info message not logged to buffer 1.\nBuffer contents: %q", buf1.String())
	}

	// Test with actual file
	tmpFile, err := os.CreateTemp("", "multilog_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Create logger with file
	logger = multilog.NewMultiWithOptions(opts, buf1, tmpFile)
	logger.Info("file test message")

	// Read file contents
	fileContents, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("File contents: %q", string(fileContents))

	if !bytes.Contains(fileContents, []byte("file test message")) {
		t.Errorf("Message not written to file.\nFile contents: %q", string(fileContents))
	}
}

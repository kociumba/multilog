
# multilog

multilog is a simple wrapper around [charm/log](https://github.com/charmbracelet/log) that enables creating loggers with multiple outputs.


## Installation

```bash
go get github.com/kociumba/multilog
```
    
## Basic usage

multilog works exatly the same as [charm/log](https://github.com/charmbracelet/log) only difference being the returned logger can write to multiple `io.Writer` interfaces simultaneously.

```go
package main

import (
    "os"

    "github.com/kociumba/multilog"
)

func main() {
    // Create a log file
    logFile, err := os.Create("log.txt")
    if err != nil {
        panic(err)
    }
    defer logFile.Close()

    // Create a new multilog logger that writes to stdout and the log file
    log := multilog.NewMulti(os.Stdout, logFile)

    log.Info("logging info into stdout and a file 😎")
}
```

> [!IMPORTANT]
> Due to the limitations of the `io.Writer` interface, which only returns a single error, multilog can only surface one error at a time when multiple writes fail. This means if multiple writers encounter errors, only the first error will be returned through [charm/log's](https://github.com/charmbracelet/log) error handling.


## Using with options

Just like in [charm/log](https://github.com/charmbracelet/log) you can create a logger with options like this:

```go
package main

import (
    "os"
    "time"

    "github.com/charmbracelet/log"
    "github.com/kociumba/multilog"
)

func main() {
    // Create a log file
    logFile, err := os.Create("log.txt")
    if err != nil {
        panic(err)
    }
    defer logFile.Close()

    // Create a new multilog logger with custom options
    log := multilog.NewMultiWithOptions(log.Options{
        ReportCaller:      true,
        ReportTimestamp:   true,
        TimeFormat:        time.RFC3339,
        Prefix:            "multilogging :3",
    }, os.Stdout, logFile)

    log.Info("logging info into stdout and a file 😎")
}
```

> [!NOTE]
> Due to Go's requirement that variadic arguments must be the last parameter, the parameter order differs from [charm/log](https://github.com/charmbracelet/log): `options, writers...` instead of `writer, options`.

The returned loggers are standard [charm/log](https://github.com/charmbracelet/log) `log.Logger` types so options can also be changed after creatation like this:

```go
log.Info("the calling function will not be reported")
log.SetReportCaller(true)
log.Info("the calling function will be reported")
```


## The multiwriter

multilog simply creates new `log.Logger` instances using the internal [multilog/multiwriter](https://github.com/kociumba/multilog/tree/main/multiwriter) package.

If you want to use this multiwriter by itself for other purposes without the [charm/log](https://github.com/charmbracelet/log) wrapper, you can use the `NewMultiWriter` function from the [multilog/multiwriter](https://github.com/kociumba/multilog/tree/main/multiwriter) package.

```go
package main

import (
    "bytes"
    "fmt"
    "net"
    "os"

    "github.com/kociumba/multilog/multiwriter"
)

func main() {
    // Create a buffer, a file writer, and a TCP writer
    buf := &bytes.Buffer{}
    logFile, err := os.Create("log.txt")
    if err != nil {
        panic(err)
    }
    defer logFile.Close()

    // Connect to a remote TCP server
    tcpConn, err := net.Dial("tcp", "remote-dashboard.example.com:12345")
    if err != nil {
        panic(err)
    }
    defer tcpConn.Close()

    // Create a new MultiWriter instance with the local and remote writers
    multi := multiwriter.NewMultiWriter(os.Stdout, buf, logFile, tcpConn)

    // Write some data to the MultiWriter
    n, err := multi.Write([]byte("Writing to multiple writers including a remote TCP server 😎"))
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Println("Wrote", n, "bytes")
}
```

You can also provide an `ErrorHandler` callback to be called for each write error that occurs.

```go
package main

import (
    "bytes"
    "fmt"
    "net"
    "os"

    "github.com/kociumba/multilog/multiwriter"
)

func main() {
    // Create a buffer, a file writer, and a TCP writer
    buf := &bytes.Buffer{}
    logFile, err := os.Create("log.txt")
    if err != nil {
        panic(err)
    }
    defer logFile.Close()

    // Connect to a remote TCP server
    tcpConn, err := net.Dial("tcp", "remote-dashboard.example.com:12345")
    if err != nil {
        panic(err)
    }
    defer tcpConn.Close()

    // Create a new MultiWriter instance with an error handler
    multi := multiwriter.NewMultiWriter(os.Stdout, buf, logFile, tcpConn, &failingWriter{})
    multi.ErrorHandler = func(we multiwriter.WriteError) {

        // You can check for error types in the ErrorHandler callback
        if we.Err == io.ErrShortWrite {
            fmt.Println("short write detected")
        }
        fmt.Printf("Write failed for writer %T: %v\n", we.Writer, we.Err)
    }

    // Write some data to the MultiWriter
    n, err := multi.Write([]byte("Writing to multiple writers with error handling 😎"))
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Println("Wrote", n, "bytes")
}

type failingWriter struct{}

func (fw *failingWriter) Write(p []byte) (int, error) {
    return 0, fmt.Errorf("this writer always fails")
}
```

> [!IMPORTANT]
> Keep in mind all of the writers that are in the multiwriter will always be called, even if there is an error in one or more. If you want to terminate the multilogger after an error you will need to handle that yourself.

## FAQ

#### Q: Why wrap [charm/log](https://github.com/charmbracelet/log) instead of creating a custom logger?

A: Simple anwser is that i really like the simplicity and design of [charm/log](https://github.com/charmbracelet/log) and this is a very niche use case that propably isn't worth making a whole new logging library for. TLDR - if not a wrapper around [charm/log](https://github.com/charmbracelet/log) then a wrapper around std/log

#### Q: What about that error handling from earlier ?

A: ~~The only half deacent idea I had for fixing this issue is to create an error channel that you as the user could listen to for errors during logging and handle them accordingly. But this brings a whole lot of complexity to this otherwise tiny wrapper so it may or may not happen at some pont.~~

This is now resolved when using the [multilog/multiwriter](https://github.com/kociumba/multilog/tree/main/multiwriter) package. You can set a custom callback function that will execute for each write error that occurs.

#### Q: Where in gods name would I need this ?

A: This isn't usefull for 99% of projects out there, but these are just some ideas I had while coding this:

- Windows services that need to log to both a file and stdout during development
- Applications that need to log locally while also sending logs to a remote system (e.g., via a TCP Writer). This includes for example servers where you could send the logs to a remote dashboard and to a local file at the same time.


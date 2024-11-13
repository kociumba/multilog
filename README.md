
# multilog

multilog is a simple wrapper around [charm/log](https://github.com/charmbracelet/log) that enables creating loggers with multiple outputs.


## Installation

```bash
  go get github.com/kociumba/multilog
```
    
## Basic usage

multilog works exatly the same as [charm/log](https://github.com/charmbracelet/log) only difference being the returned logger can write to multiple `io.Writer` interfaces simultaneously.

```go
logFile, err := os.Create("log.txt")

log := multilog.NewMulti(os.Stdout, logFile)

log.Info("logging info into stdout and a file!")
```

> [!IMPORTANT]
> Due to the limitations of the `io.Writer` interface, which only returns a single error, multilog can only surface one error at a time when multiple writes fail. This means if multiple writers encounter errors, only the first error will be returned through [charm/log's](https://github.com/charmbracelet/log) error handling.


## Using with options

Just like in [charm/log](https://github.com/charmbracelet/log) you can create a logger with options like this:

```go
logFile, err := os.Create("log.txt")

log := multilog.NewMultiWithOptions(log.Options{
    ReportCaller: true,
    ReportTimestamp: true,
    TimeFormat: time.MultiDimensional,
    Prefix: "multilogging :3",
    },
    os.Stdout,
    logFile,
)

log.Info("logging info into stdout and a file!")
```

> [!NOTE]
> Due to Go's requirement that variadic arguments must be the last parameter, the parameter order differs from [charm/log](https://github.com/charmbracelet/log): `options, writers...` instead of `writer, options`.


## FAQ

#### Q: Why wrap [charm/log](https://github.com/charmbracelet/log) instead of creating a custom logger?

A: Simple anwser is that i really like the simplicity and design of [charm/log](https://github.com/charmbracelet/log) and this is a very niche use case that propably isn't worth making a whole new logging library for. TLDR - if not a wrapper around [charm/log](https://github.com/charmbracelet/log) then a wrapper around std/log

#### Q: What about that error handling from earlier ?

A: The only half deacent idea I had for fixing this issue is to create an error channel that you as the user could listen to for errors during logging and handle them accordingly. But this brings a whole lot of complexity to this otherwise tiny wrapper so it may or may not happen at some pont. 

#### Q: Where in gods name would I need this ?

A: This isn't usefull for 99% of projects out there, but these are just some ideas I had while coding this:

- Windows services that need to log to both a file and stdout during development
- Applications that need to log locally while also sending logs to a remote system (e.g., via a TCP Writer). This includes for example servers where you could send the logs to a remote dashboard and to a local file at the same time.


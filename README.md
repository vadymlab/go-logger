# go-logger

A flexible logging library built on top of `zap` with customizable log levels, structured logging, and context
management.

## Features

- **Log Levels**: Supports standard and custom log levels including Trace, Debug, Info, Warn, Error, and Fatal.
- **Structured Logging**: Easily add key-value pairs to log messages.
- **Contextual Logging**: Attach loggers to context for efficient and structured logging throughout your application.
- **Custom Formatting**: Includes options for JSON or console output, with customizable caller information formatting.
- **Trace-Level Handling**: Provides custom handling for detailed trace logs.

## Installation

To use `go-logger` in your project, import the package:

```go
import "github.com/vadymlab/go-logger"
```

## Basic Usage

### Creating a Logger

Create a new logger with JSON or console formatting and the desired log level:

```go
config := &log.Config{
IsJson: true, // true for JSON output, false for console
Level:  "INFO", // log level: TRACE, DEBUG, INFO, WARNING, ERROR, FATAL
}

logger, err := log.NewLogger(config)
if err != nil {
panic("failed to initialize logger: " + err.Error())
}

// Set the logger as the default logger
log.SetDefaultLogger(logger)
```

### Logging Messages

You can log messages at different levels:

```go
logger.Info("This is an info message")
logger.Debugf("Debugging %s with id: %d", "test", 123)
logger.Errorw("Error occurred", "error", err)
logger.Fatal("Fatal error, application will terminate")
```

### Structured Logging

Add context to logs using key-value pairs:

```go
logger.WithField("userID", 1234).Info("User logged in")
logger.WithError(errors.New("file not found")).Error("Failed to open file")
```

### Using Context with Logger

Attach a logger to a context to maintain structured logging in applications:

```go
ctx := context.Background()
ctx = log.ToContext(ctx, logger)

// Retrieve logger from context
loggerFromCtx := log.FromContext(ctx)
loggerFromCtx.Info("Logging with context-attached logger")
```

## Advanced Usage

### Trace-Level Logging

`go-logger` provides a custom `Trace` level for detailed trace messages:

```go
logger.Trace("Detailed trace message for debugging")
```

### Skip Callers

Adjust the number of caller stack frames to skip for cleaner log outputs:

```go
logger.SkipCallers(2).Info("Adjusted caller information")
```

### Custom Formatting

Customize log formats for console output by changing the encoding and caller format:

```go
config := &log.Config{
IsJson: false, // Use console output
Level:  "DEBUG",
}

// Additional customizations for console format
logger, err := log.NewLogger(config)
```

## Log Level Conversion

Convert string log levels to `LogLevel`:

```go
logLevel := log.Text2Level("INFO")
logger := log.newZap(false, logLevel)
```

## Example Application

Below is an example of how to set up the logger and use it in a sample application:

```go
package main

import (
	"context"
	"errors"
	"log"
)

func main() {
	// Set up logger
	config := &log.Config{
		IsJson: true,
		Level:  "DEBUG",
	}

	logger, err := log.NewLogger(config)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	// Setting logger as default
	log.SetDefaultLogger(logger)

	// Example log messages
	logger.Info("Application started")
	logger.Debugf("Starting with configuration: %+v", config)

	// Using context for logger
	ctx := log.ToContext(context.Background(), logger)
	log.FromContext(ctx).Infow("Contextual log", "userID", 101)

	// Structured log with error context
	err = errors.New("sample error")
	logger.WithError(err).Error("An error occurred")

	logger.Fatal("Fatal error, exiting application")
}
```

## Contributing

Feel free to open issues or submit pull requests to improve this library!

## License

This library is open-source and available under the MIT License.

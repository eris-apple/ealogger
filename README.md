# ealogger

`ealogger` is a powerful and flexible logging library for Go that allows you to customize the logging process to suit your needs. It supports logging to files, the console, external log collection systems, and provides convenient customization tools.

## Key Features

### Flexible Logging Configuration
- **Supports logging to:**
  - File
  - Console
  - External log collectors (e.g., Graylog). Easy integration of custom adapters for log collection and processing.

### Console Log Color Customization
- Ability to set custom HEX colors for each log level.
- Ability to set custom HEX colors for message text and timestamps.

### Log File Management
- Utilizes the powerful [lumberjack](https://github.com/natefinch/lumberjack) library for log file management:
  - File size limitations.
  - Automatic archiving.
  - Convenient file rotation.

### Multiple Log Levels Supported
- Unselected
- Debug
- Info
- Warn
- Error
- Fatal

### Logging Methods with Suffixes
- **Methods with the `n` suffix** allow specifying an additional `name` field:
  ```go
  logger.Infon("TestLogger", "my log info") // output: 2024-12-13 17:21:57 INFO [TestLogger]: my log info
  ```

- **Methods with the `f` suffix** allow formatted string logging:
  ```go
  logger.Infof("%s string", "formatted") // output: 2024-12-13 17:21:57 INFO formatted string
  ```

### Additional Methods for Logging
- **DebugJSON** for formatting JSON objects:
  ```go
  logger.DebugJSON(shared.LogData{Fields: shared.LogField{"some": "field"}})
  // output: 2024-12-13 17:21:57 DEBUG {
  //   "Fields": {
  //     "some": "field"
  //   },
  //   "Error": null,
  //   "TraceName": "",
  //   "WithName": false
  // }
  ```

- **WithFields**, **WithError**, **WithName** for detailed context:
  ```go
  logger.WithFields(shared.LogField{"some": "field"}).Info("with fields") // output: 2024-12-13 17:21:57 INFO with fields some=field
  logger.WithError(errors.New("my error")).Info("with error")  // output: 2024-12-13 17:21:57 INFO with error err=my error
  logger.WithName("TestLogger").Info("with name") // output: 2024-12-13 17:21:57 INFO [TestLogger]: with name
  ```

## Creating a Custom Adapter

ealogger provides an interface for creating custom adapters to handle logs in a specific way:
```go
type Adapter interface {
    Log(log shared.Log)
    Format(log *shared.Log)
}
```

To create a custom adapter, implement this interface. For example:
```go
type TestAdapter struct {
    writer *MyWriter
    cfg    *TestConfig
}

func (a *TestAdapter) Log(log shared.Log) {
    // If needed, format/transform the log before processing
    a.Format(&log)
}

func (a *TestAdapter) Format(log *shared.Log) {}

// Example configuration
func NewTestAdapter(level shared.Level) *TestAdapter {
    return &TestAdapter{
        // Your adapter initialization logic
    }
}
```

Then, add the adapter to the logger:
```go
logger := ealogger.NewLogger(
    NewTestAdapter(shared.DebugLevel),
)
```

## Installation

Install the library using the following command:
```bash
go get github.com/eris-apple/ealogger
```

## Initialization

To get started, create a new logger and configure its adapters:
```go
logger := ealogger.NewLogger(
  adapters.NewDefaultConsoleAdapterWithLevel(shared.DebugLevel),
  adapters.NewDefaultFileAdapterWithLevel(shared.DebugLevel),
)
```

### Usage Example

```go
logger.Infof("Starting application: %s", "MyApp")
logger.WithName("InitPhase").Info("Initialization complete")
logger.WithError(errors.New("connection lost")).Error("Failed to connect to database")
logger.DebugJSON(map[string]interface{}{"key": "value"})
```

## Contributing

If you have suggestions or have found a bug, please create an issue or open a pull request in the [project repository](https://github.com/eris-apple/ealogger).

## License

This project is distributed under the MIT license. See the [LICENSE](LICENSE) file for details.


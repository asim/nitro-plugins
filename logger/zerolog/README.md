# Zerolog

[Zerolog](https://github.com/rs/zerolog) logger implementation for __go-micro__ [meta logger](https://github.com/micro/go-micro/tree/master/logger).

## Usage

```go
func ExampleWithOut() {
  logger.DefaultLogger = zerolog.NewLogger(logger.WithOut(os.Stdout), logger.WithLevel(logger.DebugLevel))

  logger.Logf(logger.InfoLevel, "testing: %s", "logf")

  // Output:
  // {"level":"info","message":"testing: logf"}
}
```

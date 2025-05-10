# http

## http/compress
Http middleware (inbound - todo) and RoundTripper (outbound) that handles compression (br.deflate,gzip,zstd).


## http/log
Http middleware (inbound) and RoundTripper (outbound) using slog.

### Inbound Logger (http middleware)

Inbound logger middleware (http handler) can be initialized with

#### A logger ([WithLogger](httplog/logger.go#L7))
The logger that is set with this function is the logger will be used to log the traffic. If no logger is set the `slog.Default()` will be used.

#### A level ([WithLogInLevel](httplog/logger.go#L11))
The level that is set with this function is the log level in which the middleware will log the traffic.
Default value: `slog.LevelDebug`.

#### A log mode ([WithMode](httplog/logger.go#L15))
The log mode can take two values `Drain` and `Tee`.
  * When `Drain` is selected the body of the income request is read entirely upon receiving and a copy of the body will be passed to the next http handlers.
  * When `Tee` is selected a tee reader wraps incoming request's body and then the request is passed to the next http handlers. The request body will be read when (and only) the next (or final) http handlers, read it.
Default value: `Drain`.

#### A log policy ([WithLogPolicy](httplog/logger.go#L19))
The log policy can indicate conditions

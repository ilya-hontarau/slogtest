# slogtest

Easy assertion of log/slog.

# Usage

``` go
		matcher := NewMatcher(t).WithMsg("test")
		defer matcher.Finish() // required

		logger := slog.New(matcher.Handler())
		logger.Info("test")
```

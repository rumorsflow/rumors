# Rumors

```shell
RUMORS_DEBUG=true \
RUMORS_LOG_COLORED=true \
RUMORS_TELEGRAM_OWNER= \
RUMORS_TELEGRAM_TOKEN= \
go run . serve
```

## Build

```shell
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=`date -u +1.0.0.%Y%m%d.%H%M%S`" .
```
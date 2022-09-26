FROM golang:1.19-alpine as build

WORKDIR /app

COPY . .

RUN apk update \
	&& apk add --no-cache build-base ca-certificates \
	&& update-ca-certificates \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=`date -u +1.0.0.%Y%m%d.%H%M%S`" -o release/ .

FROM scratch

LABEL org.opencontainers.image.authors="Igor Agapie <igoragapie@gmail.com>"
LABEL org.opencontainers.image.vendor="Rumors"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/release /

CMD ["/rumors"]
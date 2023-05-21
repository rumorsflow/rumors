FROM node:18-alpine as front-ui

RUN apk add --no-cache git
RUN git clone https://github.com/rumorsflow/ui.git /app

WORKDIR /app

RUN npm ci && VITE_APP_NAME_CAPTION=BETA npm run build

# ----------------------------------------------------------

FROM node:18-alpine as sys-ui

RUN apk add --no-cache git

RUN git clone https://github.com/rumorsflow/sys-ui.git /app

WORKDIR /app

RUN npm ci && VITE_APP_API_URL=/sys/api npm run build

# ----------------------------------------------------------

FROM golang:1.20-alpine as build

RUN apk update \
	&& apk add --no-cache build-base ca-certificates \
	&& update-ca-certificates

ARG VERSION=(untracked)

WORKDIR /app

COPY . .
COPY --from=front-ui /app/dist /app/internal/http/front/ui
COPY --from=sys-ui /app/dist /app/internal/http/sys/ui

RUN CGO_ENABLED=0 go build -trimpath -tags=sys_ui,ui -ldflags="-s -w -X main.version=${VERSION}" -o release/ .

# ----------------------------------------------------------

FROM scratch

LABEL org.opencontainers.image.authors="Igor Agapie <igoragapie@gmail.com>"
LABEL org.opencontainers.image.vendor="Rumors"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/release /
COPY --from=build /app/config.yaml /config.yaml

CMD ["/rumors"]
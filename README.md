[//]: # (<p align="center">)

[//]: # (    <a href="https://www.rumorsflow.com/" target="_blank" rel="noopener">)

[//]: # (        <img src="https://i.imgur.com/4BAQDhD.png" alt="RumorsFlow" />)

[//]: # (    </a>)

[//]: # (</p>)

<p align="center">
    <a href="https://github.com/rumorsflow/rumors/actions/workflows/release.yaml" target="_blank" rel="noopener"><img src="https://github.com/rumorsflow/rumors/actions/workflows/release.yaml/badge.svg" alt="build" /></a>
    <a href="https://github.com/rumorsflow/rumors/releases" target="_blank" rel="noopener"><img src="https://img.shields.io/github/v/release/rumorsflow/rumors.svg" alt="Latest releases" /></a>
    <a href="https://github.com/rumorsflow/rumors/blob/main/LICENSE" target="_blank" rel="noopener"><img src="https://img.shields.io/dub/l/vibe-d.svg" alt="LICENSE" /></a>
</p>

Welcome to [Rumors](https://www.rumorsflow.com/), a news aggregation application that brings together the latest news and updates from various sources.

Rumors parses RSS and sitemap XML files to gather and organize news and information, making it easier for you to stay informed and up-to-date on the latest developments in your field of interest.

### Bot commands

```shell
rumors - <index> <size> <search>
sites - List of available sites
sub - List of subscribed sites
on - subscribe <site>
off - unsubscribe <site>
```

### Generate Swagger Documentation 2.0

Install swag

```shell
go install github.com/swaggo/swag/cmd/swag@aa3e8d5fa2f6ee3a56f54c7ae3bd18145783eaac
```

Run the following command to generate Front API Documentation

```shell
swag i -g swagger.go -dir internal/http/front,internal/entity,internal/model,${HOME}/go/pkg/mod/github.com/gowool/wool@v0.0.0-20230509175958-16e9f1d12396 --instanceName front
```

Run the following command to generate System API Documentation

```shell
swag i -g swagger.go -dir internal/http/sys,internal/entity,${HOME}/go/pkg/mod/github.com/gowool/wool@v0.0.0-20230509175958-16e9f1d12396 --instanceName sys
```

### Generate a self-signed X.509 TLS certificate

Run the following command to generate `cert.pem` and `key.pem` files:

```shell
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

Run the following command to generate RSA private key

```shell
openssl genrsa -out rsa_key.pem 4096
```

### Docker Compose

```shell
# generate .env
cat > .env << "EOF"
RUMORS_TELEGRAM_TOKEN=<telegram bot token>
RUMORS_TELEGRAM_OWNER=<telegram owner ID>
RUMORS_HTTP_CERT_FILE=/absolute/path/cert.pem
RUMORS_HTTP_KEY_FILE=/absolute/path/key.pem
RUMORS_HTTP_JWT_PRIVATE_KEY=/absolute/path/rsa_key.pem
EOF

# start docker compose
docker compose up -d

# create new user and 2FA
docker compose exec rumors /rumors sys user create -u username -p password -e username@mail.com
```

## License

Distributed under MIT License, please see license file within the code for more details.

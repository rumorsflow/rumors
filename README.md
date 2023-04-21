<p align="center">
    <a href="https://www.rumorsflow.com/" target="_blank" rel="noopener">
        <img src="https://i.imgur.com/4BAQDhD.png" alt="RumorsFlow" />
    </a>
</p>

<p align="center">
    <a href="https://github.com/rumorsflow/rumors/actions/workflows/release.yaml" target="_blank" rel="noopener"><img src="https://github.com/rumorsflow/rumors/actions/workflows/release.yaml/badge.svg" alt="build" /></a>
    <a href="https://github.com/rumorsflow/rumors/releases" target="_blank" rel="noopener"><img src="https://img.shields.io/github/v/release/rumorsflow/rumors.svg" alt="Latest releases" /></a>
    <a href="https://github.com/rumorsflow/rumors/blob/main/LICENSE" target="_blank" rel="noopener"><img src="https://img.shields.io/dub/l/vibe-d.svg" alt="LICENSE" /></a>
</p>

```shell
go run . --dotenv=.env
```

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
swag i -g swagger.go -dir internal/http/front,internal/entity,internal/pubsub,${HOME}/go/pkg/mod/github.com/gowool/wool@v0.0.0-20230212000935-245e67db993b --instanceName front
```

Run the following command to generate System API Documentation

```shell
swag i -g swagger.go -dir internal/http/sys,internal/entity,${HOME}/go/pkg/mod/github.com/gowool/wool@v0.0.0-20230212000935-245e67db993b --instanceName sys
```

### Generate a self-signed X.509 TLS certificate

Run the following command to generate `cert.pem` and `key.pem` files:

```shell
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

Run the following command to generate RSA private key

```shell
openssl genrsa -out rsa_prv.pem 4096
```

## License

Distributed under MIT License, please see license file within the code for more details.

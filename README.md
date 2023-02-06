# Rumors v2

![License](https://img.shields.io/dub/l/vibe-d.svg)

```shell
go run . --dotenv=.env
```

### Bot commands

```shell
rumors - <index> <size> <search>
sources - List of available sources
sub - List of subscribed sources
on - subscribe <source>
off - unsubscribe <source>
```

### Generate Swagger Documentation 2.0

Install swag

```shell
go install github.com/swaggo/swag/cmd/swag@aa3e8d5fa2f6ee3a56f54c7ae3bd18145783eaac
```

Run the following command to generate Front API Documentation

```shell
swag i -g swagger.go -dir internal/http/front,internal/entity,internal/pubsub,${HOME}/go/pkg/mod/github.com/gowool/wool@v0.0.0-20230206095925-11fec9706d35 --instanceName front
```

Run the following command to generate System API Documentation

```shell
swag i -g swagger.go -dir internal/http/sys,internal/entity,${HOME}/go/pkg/mod/github.com/gowool/wool@v0.0.0-20230206095925-11fec9706d35 --instanceName sys
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

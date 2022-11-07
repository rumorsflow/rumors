# Rumors

![License](https://img.shields.io/dub/l/vibe-d.svg)

```shell
go run . --dotenv=.env serve
```

### Bot commands

```shell
rumors - <index> <size> <search>
sources - List of available sources
subscribed - List of subscribed sources
on - /subscribed alias
subscribe - <source>
sub - /subscribe alias
unsubscribe - <source>
unsub - /unsubscribe alias
```

### Generate a self-signed X.509 TLS certificate

Run the following command to generate `cert.pem` and `key.pem` files:

```shell
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

## License

Distributed under MIT License, please see license file within the code for more details.

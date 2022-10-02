# Rumors

```shell
RUMORS_DEBUG=true \
RUMORS_LOG_CONSOLE=true \
RUMORS_TELEGRAM_OWNER= \
RUMORS_TELEGRAM_TOKEN= \
RUMORS_MONGODB_URI= \
go run . serve
```

### Bot commands

```shell
add - <feed link> <lang>
rumors - <index> <size> <link search>
feed - [list: <index> <size> <link search>] or [info: <id>]
room - [list: <index> <size> <title search>] or [info: <id>]
```

Commands `/add` and `/rumors` are public. Commands `/feed` and `/room` are private and can be access only by bot owner.

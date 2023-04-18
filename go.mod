module github.com/rumorsflow/rumors/v2

go 1.20

replace (
	github.com/abadojack/whatlanggo v1.0.1 => github.com/retarus/whatlanggo v1.1.1
	github.com/oxffaa/gopher-parse-sitemap v0.0.0-20191021113419-005d2eb1def4 => github.com/rumorsflow/gopher-parse-sitemap v0.0.0-20230322153900-5684f39da055
)

require (
	github.com/abadojack/whatlanggo v1.0.1
	github.com/dlclark/regexp2 v1.8.1
	github.com/fatih/color v1.15.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/goccy/go-json v0.10.0
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/google/uuid v1.3.0
	github.com/gowool/middleware/cors v0.0.0-20230418213552-570f65c45a5d
	github.com/gowool/middleware/gzip v0.0.0-20230418213552-570f65c45a5d
	github.com/gowool/middleware/keyauth v0.0.0-20230418213552-570f65c45a5d
	github.com/gowool/middleware/prometheus v0.0.0-20230418213552-570f65c45a5d
	github.com/gowool/middleware/proxy v0.0.0-20230418213552-570f65c45a5d
	github.com/gowool/middleware/sse v0.0.0-20230418213552-570f65c45a5d
	github.com/gowool/swagger v0.0.0-20230206100617-5cfcdc7729c5
	github.com/gowool/wool v0.0.0-20230212000935-245e67db993b
	github.com/hibiken/asynq v0.24.0
	github.com/hibiken/asynq/x v0.0.0-20230106040302-cc777ebdaa62
	github.com/joho/godotenv v1.5.1
	github.com/mdp/qrterminal/v3 v3.0.0
	github.com/mergestat/timediff v0.0.3
	github.com/microcosm-cc/bluemonday v1.0.23
	github.com/mmcdole/gofeed v1.2.1
	github.com/otiai10/opengraph/v2 v2.1.0
	github.com/oxffaa/gopher-parse-sitemap v0.0.0-20191021113419-005d2eb1def4
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.4.0
	github.com/prometheus/client_golang v1.14.0
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.15.0
	github.com/swaggo/swag v1.8.11-0.20230125210707-aa3e8d5fa2f6
	go.mongodb.org/mongo-driver v1.11.3
	go.uber.org/multierr v1.10.0
	golang.org/x/crypto v0.5.0
	golang.org/x/exp v0.0.0-20230125214544-b3c2aaf6208d
	golang.org/x/net v0.8.0
	golang.org/x/sync v0.1.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mmcdole/goxpp v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/swaggo/files v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/qr v0.2.0 // indirect
)

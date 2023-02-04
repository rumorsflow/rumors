package mongodb

type Config struct {
	URI  string `mapstructure:"uri"`
	Ping bool   `mapstructure:"ping"`
}

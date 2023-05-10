package logger

import "golang.org/x/exp/slog"

// ChannelConfig configures loggers per channel.
type ChannelConfig struct {
	// Dedicated channels per logger. By default logger allocated via named logger.
	Channels map[string]*Config `mapstructure:"channels"`
}

type Config struct {
	// When AddSource is true, the handler adds a ("source", "file:line")
	// attribute to the output indicating the source code position of the log
	// statement. AddSource is false by default to skip the cost of computing
	// this information.
	AddSource bool `mapstructure:"add_source"`

	// Level is the minimum enabled logging level.
	Level string `mapstructure:"level"`

	// Encoding sets the logger's encoding. Init values are "json", "text" and "console"
	Encoding string `mapstructure:"encoding"`

	// Output is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `mapstructure:"output_paths"`

	Attrs map[string]any `mapstructure:"attributes"`
}

func (cfg *Config) OpenSinks() (WriteSyncer, error) {
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stderr"}
	}

	sink, _, err := Open(cfg.OutputPaths...)
	return sink, err
}

func (cfg *Config) Opts() HandlerOptions {
	return HandlerOptions{HandlerOptions: slog.HandlerOptions{
		Level:     ToLeveler(cfg.Level),
		AddSource: cfg.AddSource,
	}}
}

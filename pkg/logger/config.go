package logger

import (
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Config struct {
	// When AddSource is true, the handler adds a ("source", "file:line")
	// attribute to the output indicating the source code position of the log
	// statement. AddSource is false by default to skip the cost of computing
	// this information.
	AddSource bool `mapstructure:"add_source"`

	// Level is the minimum enabled logging level.
	Level string `mapstructure:"level"`

	// Encoding sets the logger's encoding. Init values are "json", "text" and
	// "console"
	Encoding string `mapstructure:"encoding"`

	// Output is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `mapstructure:"output_paths"`

	// File logger options
	FileLogger *FileConfig `mapstructure:"file_logger_options"`

	Attrs map[string]any `mapstructure:"attributes"`
}

func (cfg *Config) OpenSinks() (WriteSyncer, error) {
	if cfg.FileLogger == nil {
		if len(cfg.OutputPaths) == 0 {
			cfg.OutputPaths = []string{"stderr"}
		}

		sink, _, err := Open(cfg.OutputPaths...)
		return sink, err
	}

	cfg.FileLogger.Init()

	return AddSync(&lumberjack.Logger{
		Filename:   cfg.FileLogger.LogOutput,
		MaxSize:    cfg.FileLogger.MaxSize,
		MaxAge:     cfg.FileLogger.MaxAge,
		MaxBackups: cfg.FileLogger.MaxBackups,
		Compress:   cfg.FileLogger.Compress,
	}), nil
}

func (cfg *Config) Opts() HandlerOptions {
	return HandlerOptions{HandlerOptions: slog.HandlerOptions{
		Level:     ToLeveler(cfg.Level),
		AddSource: cfg.AddSource,
	}}
}

type FileConfig struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	LogOutput string `mapstructure:"log_output"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `mapstructure:"max_size"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `mapstructure:"max_age"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `mapstructure:"max_backups"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `mapstructure:"compress"`
}

func (fl *FileConfig) Init() {
	if fl.LogOutput == "" {
		fl.LogOutput = os.TempDir()
	}

	if fl.MaxSize == 0 {
		fl.MaxSize = 100
	}

	if fl.MaxAge == 0 {
		fl.MaxAge = 24
	}

	if fl.MaxBackups == 0 {
		fl.MaxBackups = 10
	}
}

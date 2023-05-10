// MIT License
//
// Copyright (c) 2022 Spiral Scout
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	OpNew          = "configurer: new ->"
	OpUnmarshalKey = "configurer: unmarshal key ->"
	OpUnmarshal    = "configurer: unmarshal ->"
	OpOverwrite    = "configurer: overwrite ->"
	OpParseFlag    = "configurer: parse flag ->"
)

type Configurer interface {
	// UnmarshalKey takes a single key and unmarshal it into a Struct.
	UnmarshalKey(name string, out any) error

	// Unmarshal the config into a Struct. Make sure that the tags
	// on the fields of the structure are properly set.
	Unmarshal(out any) error

	// Overwrite used to overwrite particular values in the unmarshalled config
	Overwrite(values map[string]any) error

	// Get used to get config section
	Get(name string) any

	// Has checks if config section exists.
	Has(name string) bool

	GracefulTimeout() time.Duration

	// Version returns current version
	Version() string
}

type Option func(*configurer)

type configurer struct {
	viper     *viper.Viper
	path      string
	prefix    string
	tp        string
	readInCfg []byte
	// user defined Flags in the form of <option>.<key> = <value>
	// which overwrites initial config key
	flags []string

	// Timeout ...
	timeout time.Duration
	version string
}

func WithPath(path string) Option {
	return func(c *configurer) {
		c.path = path
	}
}

func WithPrefix(prefix string) Option {
	return func(c *configurer) {
		c.prefix = prefix
	}
}

func WithConfigType(tp string) Option {
	return func(c *configurer) {
		c.tp = tp
	}
}

func WithReadInCfg(readInCfg []byte) Option {
	return func(c *configurer) {
		c.readInCfg = readInCfg
	}
}

func WithFlags(flags []string) Option {
	return func(c *configurer) {
		c.flags = flags
	}
}

func NewConfigurer(version string, timeout time.Duration, options ...Option) (Configurer, error) {
	c := &configurer{viper: viper.New(), timeout: timeout, version: version}

	for _, opt := range options {
		opt(c)
	}

	// If user provided []byte data with config, read it and ignore Path and Prefix
	if c.readInCfg != nil && c.tp != "" {
		c.viper.SetConfigType(c.tp)
		err := c.viper.ReadConfig(bytes.NewBuffer(c.readInCfg))
		return c, err
	}

	// read in environment variables that match
	c.viper.AutomaticEnv()
	if c.prefix == "" {
		return nil, fmt.Errorf("%s prefix should be set", OpNew)
	}

	c.viper.SetEnvPrefix(c.prefix)
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if c.path == "" {
		ex, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("%s %w", OpNew, err)
		}
		c.viper.AddConfigPath(filepath.Dir(ex))
		c.viper.AddConfigPath(filepath.Join("/", "etc", filepath.Base(ex)))
	} else {
		c.viper.SetConfigFile(c.path)
	}

	err := c.viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("%s %w", OpNew, err)
	}

	// automatically inject ENV variables using ${ENV} pattern
	for _, key := range c.viper.AllKeys() {
		val := c.viper.Get(key)
		switch t := val.(type) {
		case string:
			// for string just expand it
			c.viper.Set(key, parseEnvDefault(t))
		case []any:
			// for slice -> check if it's slice of strings
			strArr := make([]string, 0, len(t))
			for i := 0; i < len(t); i++ {
				if valStr, ok := t[i].(string); ok {
					strArr = append(strArr, parseEnvDefault(valStr))
					continue
				}

				c.viper.Set(key, val)
			}

			// we should set the whole array
			if len(strArr) > 0 {
				c.viper.Set(key, strArr)
			}
		default:
			c.viper.Set(key, val)
		}
	}

	// override config flags
	if len(c.flags) > 0 {
		for _, f := range c.flags {
			key, val, errP := parseFlag(f)
			if errP != nil {
				return nil, fmt.Errorf("%s %w", OpNew, errP)
			}
			c.viper.Set(key, parseEnvDefault(val))
		}
	}

	return c, nil
}

func (cfg *configurer) UnmarshalKey(name string, out any) error {
	if err := cfg.viper.UnmarshalKey(name, out); err != nil {
		return fmt.Errorf("%s %w", OpUnmarshalKey, err)
	}
	return nil
}

func (cfg *configurer) Unmarshal(out any) error {
	if err := cfg.viper.Unmarshal(out); err != nil {
		return fmt.Errorf("%s %w", OpUnmarshal, err)
	}
	return nil
}

func (cfg *configurer) Overwrite(values map[string]any) error {
	for key, value := range values {
		cfg.viper.Set(key, value)
	}
	return nil
}

func (cfg *configurer) Get(name string) any {
	return cfg.viper.Get(name)
}

func (cfg *configurer) Has(name string) bool {
	return cfg.viper.IsSet(name)
}

func (cfg *configurer) Version() string {
	return cfg.version
}

func (cfg *configurer) GracefulTimeout() time.Duration {
	return cfg.timeout
}

func parseFlag(flag string) (string, string, error) {
	if !strings.Contains(flag, "=") {
		return "", "", fmt.Errorf("%s invalid flag `%s`", OpParseFlag, flag)
	}

	parts := strings.SplitN(strings.TrimLeft(flag, " \"'`"), "=", 2)
	if len(parts) < 2 {
		return "", "", errors.New("usage: -o key=value")
	}

	if parts[0] == "" {
		return "", "", errors.New("key should not be empty")
	}

	if parts[1] == "" {
		return "", "", errors.New("value should not be empty")
	}

	return strings.Trim(parts[0], " \n\t"), parseValue(strings.Trim(parts[1], " \n\t")), nil
}

func parseValue(value string) string {
	escape := []rune(value)[0]

	if escape == '"' || escape == '\'' || escape == '`' {
		value = strings.Trim(value, string(escape))
		value = strings.ReplaceAll(value, fmt.Sprintf("\\%s", string(escape)), string(escape))
	}

	return value
}

func parseEnvDefault(val string) string {
	// tcp://127.0.0.1:${RPC_PORT:-36643}
	// for envs like this, part would be tcp://127.0.0.1:
	return ExpandVal(val, os.Getenv)
}

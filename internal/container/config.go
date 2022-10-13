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

package container

import (
	"github.com/rumorsflow/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type Config struct {
	GracePeriod time.Duration
	PrintGraph  bool
	Logger      *zap.Logger
}

const (
	endureKey          = "endure"
	logsKey            = "logs"
	defaultGracePeriod = 30 * time.Second
)

// NewConfig creates endure container configuration.
func NewConfig(cfgFile, prefix string) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.SetEnvPrefix(prefix)
	v.SetConfigFile(cfgFile)

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfgLogger := logger.Config{
		Mode:  "production",
		Level: "error",
	}

	cfg := Config{
		GracePeriod: defaultGracePeriod,
		PrintGraph:  false,
	}

	if !v.IsSet(endureKey) {
		cfg.Logger, err = cfgLogger.BuildLogger()
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}

	for _, key := range v.AllKeys() {
		if !strings.HasPrefix(key, endureKey) && !strings.HasPrefix(key, logsKey) {
			continue
		}
		val := v.Get(key)
		switch t := val.(type) {
		case string:
			v.Set(key, os.ExpandEnv(t))
		case []any:
			strArr := make([]string, 0, len(t))
			for i := 0; i < len(t); i++ {
				if valStr, ok := t[i].(string); ok {
					strArr = append(strArr, os.ExpandEnv(valStr))
					continue
				}
				v.Set(key, val)
			}
			if len(strArr) > 0 {
				v.Set(key, strArr)
			}
		default:
			v.Set(key, val)
		}
	}

	cfgEndure := struct {
		GracePeriod time.Duration `mapstructure:"grace_period"`
		PrintGraph  bool          `mapstructure:"print_graph"`
	}{}

	err = v.UnmarshalKey(endureKey, &cfgEndure)
	if err != nil {
		return nil, err
	}

	if v.IsSet(logsKey) {
		err = v.UnmarshalKey(logsKey, &cfgLogger)
		if err != nil {
			return nil, err
		}
	}

	cfg.Logger, err = cfgLogger.BuildLogger()
	if err != nil {
		return nil, err
	}

	cfg.PrintGraph = cfgEndure.PrintGraph
	if cfgEndure.GracePeriod != 0 {
		cfg.GracePeriod = cfgEndure.GracePeriod
	}

	return &cfg, nil
}

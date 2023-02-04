package config

import (
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	OpUnmarshalKey errs.Op = "configurer: unmarshal key"
	OpUnmarshal    errs.Op = "configurer: unmarshal"
	OpOverwrite    errs.Op = "configurer: overwrite"
)

func UnmarshalKeyE[T any](cfg Configurer, key string) (value T, err error) {
	if !cfg.Has(key) {
		return value, errs.Errorf(OpUnmarshalKey, "config key `%s` is required", key)
	}

	err = cfg.UnmarshalKey(key, &value)
	return
}

func UnmarshalKey[T any](cfg Configurer, key string) (value T, err error) {
	if cfg.Has(key) {
		err = cfg.UnmarshalKey(key, &value)
		return
	}

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		value = reflect.New(v.Type().Elem()).Interface().(T)
	}
	return
}

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
}

type configurer struct {
	viper *viper.Viper
}

func NewConfigurer(cfgFile, prefix string) Configurer {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.SetEnvPrefix(prefix)

	if cfgFile == "" {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		v.AddConfigPath(filepath.Dir(ex))
		v.AddConfigPath(filepath.Join("/", "etc", filepath.Base(ex)))
	} else {
		v.SetConfigFile(cfgFile)
	}

	_ = v.ReadInConfig()

	for _, key := range v.AllKeys() {
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

	return &configurer{viper: v}
}

func (cfg *configurer) UnmarshalKey(name string, out any) error {
	if err := cfg.viper.UnmarshalKey(name, out); err != nil {
		return errs.E(OpUnmarshalKey, err)
	}
	return nil
}

func (cfg *configurer) Unmarshal(out any) error {
	if err := cfg.viper.Unmarshal(out); err != nil {
		return errs.E(OpUnmarshal, err)
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

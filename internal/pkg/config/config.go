package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var configFileNotFoundError viper.ConfigFileNotFoundError

// Load config file from path and bind values to cfg
//
// err := config.Load("config.yml", &cfg)
//
//	if err != nil {
//	    ...
//	}
func Load(path string, cfg any) error {
	viper.SetConfigFile(path)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil && !errors.As(err, &configFileNotFoundError) {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	bindValues(cfg)

	err = viper.Unmarshal(cfg)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	return nil
}

// FilePath get config file path
func FilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd error: %w", err)
	}

	configFile := "config.yml"
	env := strings.ToLower(os.Getenv("ENV"))
	if env != "" {
		configFile = fmt.Sprintf("config.%s.yml", env)
	}

	path := filepath.ToSlash(fmt.Sprintf("%s/%s", dir, configFile))
	return path, nil
}

func bindValues(iface any, parts ...string) {
	ift := reflect.TypeOf(iface)
	if ift != nil && ift.Kind() == reflect.Pointer {
		ift = ift.Elem()
	}

	ifv := reflect.ValueOf(iface)
	if ifv.Kind() == reflect.Pointer {
		ifv = ifv.Elem()
	}

	processField(ifv, ift, parts)
}

func processField(v reflect.Value, t reflect.Type, parts []string) {
	for i := 0; i < t.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		tag, ok := fieldType.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		tags := append(parts, tag)
		if fieldVal.Kind() == reflect.Struct {
			processField(fieldVal, fieldType.Type, tags)
			continue
		}

		key := strings.Join(tags, ".")
		envKey := strings.ToUpper(strings.Join(tags, "_"))
		_ = viper.BindEnv(key, envKey)

		value, defaultValue := fieldType.Tag.Lookup("defaultvalue")
		if defaultValue {
			viper.SetDefault(key, value)
		}
	}
}

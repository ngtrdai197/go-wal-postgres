package config

import (
	"fmt"
	"github.com/inhies/go-bytesize"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

type Schema struct {
	APIInfo struct {
		PublicApiListen string `mapstructure:"public_api_listen"`
	} `mapstructure:"api_info"`
	Database Database `mapstructure:"database"`
	Kafka    Kafka    `mapstructure:"kafka"`
}

// StringToByteSizeHookFunc returns a DecodeHookFunc that converts
// hex string to bytesize.ByteSize.
func StringToByteSizeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(bytesize.B) {
			return data, nil
		}

		sDec, err := bytesize.Parse(data.(string))
		if err != nil {
			return nil, err
		}

		return sDec, nil
	}
}

var Config *Schema

func Init() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")          // Look for config in current directory
	config.AddConfigPath("config/")    // Optionally look for config in the working directory.
	config.AddConfigPath("../config/") // Look for config needed for tests.
	config.AddConfigPath("../")        // Look for config needed for tests.

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	err = config.Unmarshal(&Config, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.RecursiveStructToMapHookFunc(),
			StringToByteSizeHookFunc(),
		),
	))
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

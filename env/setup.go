package env

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	var fileConfig string
	switch os.Getenv("APP_ENV") {
	case "DEV":
		fileConfig = "dev-env.json"
	case "UAT":
		fileConfig = "uat-env.json"
	case "PROD":
		fileConfig = "prod-env.json"
	default:
		fileConfig = "env.json"
	}

	logrus.Info("APP ENV : " + os.Getenv("APP_ENV"))
	viper.SetConfigFile(fileConfig)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	logrus.Info("Config file : " + fileConfig)
}

func String(key string, defaultValue string) string {
	value, isOK := viper.Get(key).(string)
	if !isOK {
		return defaultValue
	}
	return value
}

func Int(key string, defaultValue int) int {
	if value, isOK := viper.Get(key).(string); isOK {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return intValue
	}

	return defaultValue
}

func Bool(key string, defaultValue bool) bool {
	if value, isOK := viper.Get(key).(string); isOK {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return boolValue
	}

	return defaultValue
}

func Interface(key string, defaultValue interface{}) interface{} {
	value := viper.Get(key)
	if value == nil {
		return defaultValue
	}
	return value
}

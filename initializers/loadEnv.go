package initializers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	AccessTokenKey       string        `mapstructure:"TOKEN_KEY"`
	AccessTokenExpiresIn time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`

	KafkaBroker1 string `mapstructure:"KAFKA_PRODUCER_BROKER_1"`
	KafkaBroker2 string `mapstructure:"KAFKA_PRODUCER_BROKER_2"`
	KafkaBroker3 string `mapstructure:"KAFKA_PRODUCER_BROKER_3"`
}

var lock = &sync.Mutex{}

var config *Config

func GetConfig() *Config {
	if config == nil {
		lock.Lock()
		defer lock.Unlock()
		if config == nil {
			config = &Config{}
			err := LoadConfig(".")
			if err != nil {
				log.Fatal("? Could not load environment variables", err)
			}
		} else {
			fmt.Println("Config loaded.")
		}
	} else {
		fmt.Println("Config loaded.")
	}

	return config
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(config)
	return err
}

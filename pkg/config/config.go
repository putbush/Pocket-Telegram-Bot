package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	TelegramBotURL    string `mapstructure:"bot_url"`
	DBPath            string `mapstructure:"db_file"`

	Message Messages
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuth       string `mapstructure:"already_authorized"`
	SavesSuccessfully string `mapstructure:"saves_successfully"`
}

type Errors struct {
	Default        string `mapstructure:"default"`
	InvalidURL     string `mapstructure:"invalid_url"`
	Unauthorized   string `mapstructure:"unauthorized"`
	UnableToSave   string `mapstructure:"unable_to_save"`
	UnknownCommand string `mapstructure:"unknown_command"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Message.Errors); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.responses", &cfg.Message.Responses); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("token"); err != nil {
		return err
	}

	if err := viper.BindEnv("consumer_key"); err != nil {
		return err
	}

	if err := viper.BindEnv("auth_server_url"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("token")
	cfg.PocketConsumerKey = viper.GetString("consumer_key")
	cfg.AuthServerURL = viper.GetString("auth_server_url")

	return nil
}

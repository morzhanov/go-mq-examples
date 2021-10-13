package config

import "github.com/spf13/viper"

type Config struct {
	KafkaURI      string `mapstructure:"KAFKA_URI"`
	KafkaTopic    string `mapstructure:"KAFKA_TOPIC"`
	ActiveMQURI   string `mapstructure:"ACTIVE_MQ_URI"`
	ActiveMQQueue string `mapstructure:"ACTIVE_MQ_QUEUE"`
	RabbitMQURI   string `mapstructure:"RABBIT_MQ_URI"`
	RabbitMQQueue string `mapstructure:"RABBIT_MQ_QUEUE"`
}

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

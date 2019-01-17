package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := NewBot().Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func init() {
	pflag.StringP("server-base-url", "b", "http://localhost:8080/", "Базовый урл для коротких ссылок")
	pflag.StringP("server-host", "a", "127.0.0.1", "Адрес для запуска сервера")
	pflag.StringP("server-port", "p", ":8080", "Порт для запуска сервера")
	pflag.StringP("file-storage-path", "f", "", "Путь к файлу с ссылками")
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)

	if err != nil {
		log.Fatalln(err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

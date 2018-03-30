package main

import (
	"fmt"

	"github.com/Slava12/Computer_Market/config"
	"github.com/Slava12/Computer_Market/logger"
)

func main() {

	logger.Init("log.txt")
	configFile, errorLoadConfig := config.Parse()
	if errorLoadConfig != nil {
		logger.Fatal(errorLoadConfig, "Файл конфигурации не был загружен!")
	} else {
		logger.Info("Файл конфигурации успешно загружен.")
	}
	fmt.Println(configFile)
}

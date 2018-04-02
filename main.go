package main

import (
	"fmt"

	"github.com/Slava12/Computer_Market/config"
	"github.com/Slava12/Computer_Market/database"
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
	errorDatabase := database.Connect(configFile)
	if errorDatabase != nil {
		logger.Fatal(errorDatabase, "Не удалось подключиться к базе данных!")
	} else {
		logger.Info("Подключение к базе данных прошло успешно.")
	}
	user := database.User{0, 0, "lol", "123", "", ""}
	id, _ := database.NewUser(user.AccessLevel, user.Login, user.Password, user.Email, user.FullName)
	fmt.Println(id)
}

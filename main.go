package main

import (
	"github.com/Slava12/Computer_Market/config"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/handlefunc"
	"github.com/Slava12/Computer_Market/logger"
)

func main() {

	logger.Init("log.txt")
	configFile, errorLoadConfig := config.Parse()
	if errorLoadConfig != nil {
		logger.Error(errorLoadConfig, "Файл конфигурации не был загружен!")
	} else {
		logger.Info("Файл конфигурации успешно загружен.")
	}
	errorDatabase := database.Connect(configFile)
	if errorDatabase != nil {
		logger.Error(errorDatabase, "Не удалось подключиться к базе данных!")
	} else {
		logger.Info("Подключение к базе данных прошло успешно.")
	}

	handlefunc.InitHTTP(configFile)

	//user := database.User{0, 0, "lol", "123", "", ""}
	//_ = database.NewUser(user.AccessLevel, user.Login, user.Password, user.Email, user.FullName)
	//_ = database.DelUser(3)
	/*users, _ := database.GetAllUsers()
	fmt.Println(users)*/
}

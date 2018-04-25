package main

import (
	"github.com/Slava12/Computer_Market/config"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/files"
	"github.com/Slava12/Computer_Market/handlefunc"
	"github.com/Slava12/Computer_Market/logger"
)

func main() {
	configFile, errorLoadConfig := config.Parse()
	logsFolder := configFile.Logs.Folder
	logsName := configFile.Logs.Name
	files.CreateDirectory(logsFolder)
	logs := logsFolder + logsName
	logger.Init(logs)
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

	folder := configFile.Files.Folder
	files.CreateDirectory(folder)

	handlefunc.InitHTTP(configFile)

	//user := database.User{0, 0, "lol", "123", "", ""}
	//_ = database.NewUser(user.AccessLevel, user.Login, user.Password, user.Email, user.FullName)
	//_ = database.DelUser(3)
	/*users, _ := database.GetAllUsers()
	fmt.Println(users)*/
}

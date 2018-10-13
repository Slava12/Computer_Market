package main

import (
	"github.com/Slava12/Computer_Market/config"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/files"
	"github.com/Slava12/Computer_Market/handlefunc"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/Slava12/Computer_Market/post"
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
	}
	logger.Info("Файл конфигурации успешно загружен.")
	errorDatabase := database.Connect(configFile)
	if errorDatabase != nil {
		logger.Error(errorDatabase, "Не удалось подключиться к базе данных!")
	}
	logger.Info("Подключение к базе данных прошло успешно.")

	folder := configFile.Files.Folder
	files.CreateDirectory(folder)

	post.Init(configFile)

	handlefunc.InitRandomizer()

	handlefunc.InitHTTP(configFile)
}

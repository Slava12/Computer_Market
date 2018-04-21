package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Init logger. Input in file.
func Init(path string) {
	log.Out = os.Stdout

	filelog, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	} else {
		log.Out = filelog
	}
	Info("//---------------------Запуск сервера----------------------------//")
}

//Debug inform about bug
func Debug(message ...interface{}) {
	fmt.Println(message...)
	log.Debug(message...)
}

//Info on program execution
func Info(message ...interface{}) {
	fmt.Println(message...)
	log.Info(message...)
}

//Warn about errors
func Warn(message ...interface{}) {
	fmt.Println(message...)
	log.Warn(message...)
}

//Error program
func Error(message ...interface{}) {
	fmt.Println(message...)
	log.Error(message...)
	os.Exit(1)
}

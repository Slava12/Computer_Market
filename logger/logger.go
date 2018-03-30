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
func Debug(err error, message string) {
	fmt.Println(message)
	log.Debug(err, message)
}

//Info on program execution
func Info(message string) {
	fmt.Println(message)
	log.Info(message)
}

//Warn about errors
func Warn(err error, message string) {
	fmt.Println(message)
	log.Warn(err, message)
}

//Error program
func Error(err error, message string) {
	fmt.Println(message)
	log.Error(err, message)
	os.Exit(1)
}

//Fatal bug in program
func Fatal(err error, message string) {
	fmt.Println(message)
	log.Fatal(err, message)
	panic(message)
}

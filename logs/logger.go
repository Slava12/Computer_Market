package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Init logger. Input in file.
func Init(path string) {
	log.Out = os.Stdout

	filelog, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	} else {
		log.Out = filelog
	}
}

//Debug inform about bug
func Debug(err error, message string) {
	log.Debug(err, message)
}

//Info on program execution
func Info(message string) {
	log.Info(message)
}

//Warn about errors
func Warn(err error, message string) {
	log.Warn(err, message)
}

//Error program
func Error(err error, message string) {
	log.Error(err, message)
	os.Exit(1)
}

//Fatal bug in program
func Fatal(err error, message string) {
	log.Fatal(err, message)
	panic(message)
}

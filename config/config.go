package config

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config служит типом данных, полученных из файла конфигурации
type Config struct {
	HTTP     HTTP     `yaml:"http"`
	Database Database `yaml:"database"`
	Logs     Logs     `yaml:"logs"`
	Files    Files    `yaml:"files"`
	Post     Post     `yaml:"post"`
}

// HTTP содержит данные конфига, связанные с протоколом HTTP
type HTTP struct {
	Port string `yaml:"port"`
}

// Database содержит данные конфига, связанные с базой данных
type Database struct {
	Driver   string `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	SSLmode  string `yaml:"sslmode"`
}

// Logs содержит данные конфига, связанные с логированием
type Logs struct {
	Folder string `yaml:"folder"`
	Name   string `yaml:"name"`
}

// Files содержит данные конфига, связанные с хранением файлов на сервере
type Files struct {
	Folder string `yaml:"folder"`
}

// Post содержит данные конфига, связанные с отправкой почты
type Post struct {
	Server   string `yaml:"server"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Parse обрабатывает файл конфигурации
func Parse() (Config, error) {
	config := Config{}

	configurationPath := flag.String("path", "conf.yaml", "Путь до файла конфигурации.")
	flag.Parse()

	bytesFile, errorReadFile := ioutil.ReadFile(*configurationPath)
	if errorReadFile != nil {
		return config, errorReadFile
	}

	errorUnmarshal := yaml.Unmarshal(bytesFile, &config)
	if errorUnmarshal != nil {
		return config, errorUnmarshal
	}

	return config, nil
}

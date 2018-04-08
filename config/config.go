package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config служит типом данных, полученных из файла конфигурации
type Config struct {
	HTTP     HTTP     `yaml:"http"`
	Database Database `yaml:"database"`
	Logs     Logs     `yaml:"logs"`
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
}

// Parse обрабатывает файл конфигурации
func Parse() (Config, error) {
	config := Config{}

	//configurationPath := flag.String("path", "", "Путь до файла конфигурации.")
	configurationPath := "config.yaml" // убрать при работе
	//flag.Parse()

	bytesFile, errorReadFile := ioutil.ReadFile(configurationPath)
	if errorReadFile != nil {
		return config, errorReadFile
	}

	errorUnmarshal := yaml.Unmarshal(bytesFile, &config)
	if errorUnmarshal != nil {
		return config, errorUnmarshal
	}

	return config, nil
}

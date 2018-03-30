package config

import (
	"flag"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Http     HttpInfo     `yaml:http`
	Database DatabaseInfo `yaml:database`
	Logs     LogsInfo     `yaml:logs`
}

type HttpInfo struct {
	Port string `yaml:port`
}

type DatabaseInfo struct {
	Driver   string `yaml:driver`
	User     string `yaml:user`
	Password string `yaml:password`
	DBname   string `yaml:dbname`
	Host     string `yaml:host`
	Port     string `yaml:port`
	SSLmode  string `yaml:sslmode`
}

type LogsInfo struct {
	Folder string `yaml:folder`
}

func Parse() (Config, error) {
	config := Config{}

	configurationPath := flag.String("path", "", "Путь до файла конфигурации.")
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

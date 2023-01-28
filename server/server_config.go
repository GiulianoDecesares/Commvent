package server

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	ClientConfig ClientConfig `yaml:"client"`
}

func ConfigFromYamlFile(filePath string, fileName string) (*ServerConfig, error) {
	var err error = nil
	var serverConfig *ServerConfig = nil
	var configFile *os.File

	fullPath := path.Join(filePath, fileName)

	if configFile, err = os.Open(fullPath); err == nil {
		defer configFile.Close()

		decoder := yaml.NewDecoder(configFile)

		err = decoder.Decode(serverConfig)
	}

	return serverConfig, err
}

package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func GetConfig(data []byte) *Config {
	var config Config
	dataReader := bytes.NewReader(data)
	err := json.NewDecoder(dataReader).Decode(&config)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error on decoding json config: %v", err))
	}

	return &config
}

func GetConfigFromFile(fullPathToFile string) *Config {
	fileData, err := ioutil.ReadFile(fullPathToFile)
	if err != nil {
		log.Fatal(err)
	}

	c := GetConfig(fileData)
	return c
}

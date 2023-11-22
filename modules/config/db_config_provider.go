package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tochamateusz/micro/modules/database"
)

type DbFileConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	Port     uint   `json:"port"`
}

func ProvideDbConfig(*database.Config) (*database.Config, error) {

	//TODO: maybe later path can be provided from env var
	jsonDbConfigFile, err := os.Open("./configs/db.json")
	byteValue, _ := io.ReadAll(jsonDbConfigFile)

	if err != nil {
		return nil, &JsonConfigNotExit{}
	}

	var dbConfig DbFileConfig

	err = json.Unmarshal(byteValue, &dbConfig)
	if err != nil {
		return nil, &CannotParseConfigFile{}
	}

	config := &database.Config{
		Host:     dbConfig.Host,
		User:     dbConfig.User,
		Password: dbConfig.Password,
		DbName:   dbConfig.DbName,
		Port:     dbConfig.Port,
	}

	fmt.Printf("config: %+v\n", config)
	return config, nil

}

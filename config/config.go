package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Cfg struct {
	AuthServer struct {
		Port struct {
			Public  int `yaml:"public"`
			Private int `yaml:"private"`
		} `yaml:"port"`
	} `yaml:"auth_server"`
	UserServer struct {
		Port int `yaml:"port"`
	} `yaml:"user_server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Secret struct {
		AccessToken  string `yaml:"access_token"`
		RefreshToken string `yaml:"refresh_token"`
	} `yaml:"secret"`
	AuthServicePrivateUrl string `yaml:"auth_service_private_url"`
}

var Config Cfg

func init() {
	data, err := os.ReadFile("./config/config.yml")
	if err != nil {
		log.Printf("failed to read configuration file: %v", err)
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Printf("failed to unmarshal configuration data: %v", err)
	}

	if os.Getenv("DB_HOST") != "" {
		Config.Database.Host = os.Getenv("DB_HOST")
	}
}

package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

const (
	configFilePath = "config.yaml"
	pgConfigStr    = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
)

type PGConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
}

type APIConfig struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}

type Config struct {
	API APIConfig `yaml:"api"`
	DB 	PGConfig  `yaml:"postgres"`
}

func GetConfig() (*Config, error) {
	dbConfFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = yaml.Unmarshal(dbConfFile, &conf)
	return conf, err
}

func (c *PGConfig) ConnStr() string {
	return fmt.Sprintf(pgConfigStr,
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
	)
}

func (c *APIConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

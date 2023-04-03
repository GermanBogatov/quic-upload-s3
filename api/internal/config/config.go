package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	Namespace = "http3_upload"
)

type Config struct {
	Server Server      `yaml:"server"`
	S3     S3          `yaml:"s3"`
	Cdn    string      `yaml:"cdn"`
	Health HealthCheck `yaml:"health"`
}

type HealthCheck struct {
	CheckIntervalSec int `yaml:"check_interval" default:"60"`
}

type S3 struct {
	Host      string `yaml:"host"`
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
}

type Server struct {
	Addr            string `yaml:"addr" default:"localhost:8080"`
	ShutdownTimeout int    `yaml:"shutdown_timeout" default:"1"`
	SizeUpload      int64  `yaml:"size_upload"`
	ApiKey          string `yaml:"api_key"`
	DeleteKey       string `yaml:"delete_key"`
}

func NewConfig(configFile string) (*Config, error) {
	var config Config
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failure upload yaml file. err %w", err)
	}

	err = defaults.Set(&config)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &config, nil
}

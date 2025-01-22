package tests

import (
	"testing"
	"yaml"
)

type Config struct {
	Version     string    `yaml:"version"`
	Environment string    `yaml:"environment"`
	App         AppConfig `yaml:"app"`
	Database    Database  `yaml:"database"`
	Logging     Logging   `yaml:"logging"`
	Features    Features  `yaml:"features"`
}

type AppConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Port        int    `yaml:"port"`
	Debug       bool   `yaml:"debug"`
}

type Database struct {
	Host     string   `yaml:"host"`
	Port     int      `yaml:"port"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Schema   string   `yaml:"schema"`
	Campo    []string `yaml:"campo"`
	Retry    Retry    `yaml:"retry"`
}

type Retry struct {
	Attempts int `yaml:"attempts"`
	Delay    int `yaml:"delay"`
}

type Logging struct {
	Level    string      `yaml:"level"`
	File     string      `yaml:"file"`
	Rotation LogRotation `yaml:"rotation"`
}

type LogRotation struct {
	MaxSize  string `yaml:"max_size"`
	MaxFiles int    `yaml:"max_files"`
}

type Features struct {
	Authentication bool     `yaml:"authentication"`
	Caching        Caching  `yaml:"caching"`
	Analytics      bool     `yaml:"analytics"`
	Metadata       Metadata `yaml:"metadata"`
}

type Caching struct {
	Enabled bool   `yaml:"enabled"`
	Type    string `yaml:"type"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}

type Metadata struct {
	CreatedBy string   `yaml:"created_by"`
	CreatedAt string   `yaml:"created_at"`
	Tags      []string `yaml:"tags"`
}

func TestLoadFileYaml(t *testing.T) {

	contentYaml := Config{}

	err := yaml.Read(&contentYaml, "./test.yaml")

	if err != nil {
		t.Error(err)
	}

	//fmt.Println(contentYaml)
}

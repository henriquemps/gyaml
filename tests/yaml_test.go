package tests

import (
	"fmt"
	"testing"
	"yaml"
)

type Root struct {
	Config Config `yaml:"config"`
	Dados  Dados  `yaml:"dados"`
}

type Dados struct {
	RG     string          `yaml:"rg"`
	Doc    string          `yaml:"doc"`
	Outros OutrosEstrutura `yaml:"outros"`
}

type OutrosEstrutura struct {
	A int          `yaml:"A"`
	B int          `yaml:"B"`
	C SubEstrutura `yaml:"C"`
}

type SubEstrutura struct {
	D int `yaml:"D"`
	E int `yaml:"E"`
}

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

func TestFUnmarshal(t *testing.T) {

	contentYaml := Dados{}

	yaml.FUnmarshal(&contentYaml, "./test.yaml")
}

func TestUnmarshal(t *testing.T) {

	contentYaml := Dados{}

	yaml.Unmarshal(&contentYaml, `dados:
  rg: 000
  doc: 000
  outros:
    A: 123
    B: 123
    C:
      D: 333
      E: 444`)

	fmt.Println(contentYaml)
}

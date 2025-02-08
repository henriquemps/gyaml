package gyaml

import (
	"os"
	"testing"
)

type Root struct {
	Config Config `yaml:"config"`
}

type Config struct {
	Version     string    `yaml:"version"`
	Environment string    `yaml:"environment"`
	MultiA      string    `yaml:"multi"`
	MultiB      string    `yaml:"multi"`
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

// Feature tests

func TestFUnmarshal(t *testing.T) {

	contentYaml := Config{}

	FUnmarshal(&contentYaml, "./mocks/test.yaml")

	if contentYaml.Features.Caching.Host != "cache.example.com" {
		t.Errorf("TestFUnmarshal() error to unmarshal file")
	}
}

func TestFUnmarshalWithRoot(t *testing.T) {

	contentYaml := Root{}

	FUnmarshal(&contentYaml, "./mocks/test.yaml")

	if contentYaml.Config.Features.Caching.Host != "cache.example.com" {
		t.Errorf("TestFUnmarshalWithRoot() error to unmarshal file")
	}
}

func TestUnmarshalWithFileContent(t *testing.T) {

	contentYaml := Config{}

	content, err := os.ReadFile("./mocks/test.yaml")

	if err != nil {
		t.Fatal(err)
	}

	Unmarshal(&contentYaml, string(content))

	if contentYaml.Features.Caching.Host != "cache.example.com" {
		t.Errorf("TestUnmarshalWithFileContent() error to unmarshal file content")
	}
}

func TestUnmarshalWithFileContentWithRoot(t *testing.T) {

	contentYaml := Root{}

	content, err := os.ReadFile("./mocks/test.yaml")

	if err != nil {
		t.Fatal(err)
	}

	Unmarshal(&contentYaml, string(content))

	if contentYaml.Config.Features.Caching.Host != "cache.example.com" {
		t.Errorf("TestUnmarshalWithFileContentWithRoot() error to unmarshal file content")
	}
}

func TestUnmarshalWithRawString(t *testing.T) {

	contentYaml := Config{}

	Unmarshal(&contentYaml, `
config:
  version: 1.0.0
  environment: production
  app:
    name: SampleApp
    description: This is a sample application configuration
    port: 8080
    debug: true
  database:
    host: db.example.com
    port: 5432
    campo:
      - A
      - B
      - C
    username: db_user
    password: db_pass
    schema: public
    retry:
      attempts: 5
      delay: 2000
  logging:
    level: info
    file: /var/log/sampleapp.log
    rotation:
      max_size: 10MB
      max_files: 5
  features:
    authentication: true
    caching:
      enabled: true
      type: redis
      host: cache.example.com
      port: 6379
    analytics: false
    metadata:
      created_by: admin
      created_at: 2025-01-11
      tags:
        - app
        - config
        - yaml`)

	if contentYaml.Features.Caching.Host != "cache.example.com" {
		t.Errorf("UnmarshalWithFileContent() error to unmarshal string")
	}
}

func TestUnmarshalWithRawStringWithRoot(t *testing.T) {

	contentYaml := Root{}

	Unmarshal(&contentYaml, `
config:
  version: 1.0.0
  environment: production
  app:
    name: SampleApp
    description: This is a sample application configuration
    port: 8080
    debug: true
  database:
    host: db.example.com
    port: 5432
    campo:
      - A
      - B
      - C
    username: db_user
    password: db_pass
    schema: public
    retry:
      attempts: 5
      delay: 2000
  logging:
    level: info
    file: /var/log/sampleapp.log
    rotation:
      max_size: 10MB
      max_files: 5
  features:
    authentication: true
    caching:
      enabled: true
      type: redis
      host: cache.example.com
      port: 6379
    analytics: false
    metadata:
      created_by: admin
      created_at: 2025-01-11
      tags:
        - app
        - config
        - yaml`)

	if contentYaml.Config.Features.Caching.Host != "cache.example.com" {
		t.Errorf("UnmarshalWithFileContent() error to unmarshal string")
	}
}

// Benchmark tests

func BenchmarkFUnmarshal(b *testing.B) {

	for i := 0; i < b.N; i++ {
		contentYaml := Config{}

		FUnmarshal(&contentYaml, "./mocks/test.yaml")

		if contentYaml.Features.Caching.Host != "cache.example.com" {
			b.Errorf("TestFUnmarshal() error to unmarshal file")
		}
	}
}

func BenchmarkFUnmarshalWithRoot(b *testing.B) {

	for i := 0; i < b.N; i++ {
		contentYaml := Root{}

		FUnmarshal(&contentYaml, "./mocks/test.yaml")

		if contentYaml.Config.Features.Caching.Host != "cache.example.com" {
			b.Errorf("TestFUnmarshalWithRoot() error to unmarshal file")
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {

	for i := 0; i < b.N; i++ {
		contentYaml := Config{}

		content, err := os.ReadFile("./mocks/test.yaml")

		if err != nil {
			b.Fatal(err)
		}

		Unmarshal(&contentYaml, string(content))

		if contentYaml.Features.Caching.Host != "cache.example.com" {
			b.Errorf("BenchmarkUnmarshal() error to unmarshal file content")
		}
	}
}

func BenchmarkUnmarshalWithRoot(b *testing.B) {

	for i := 0; i < b.N; i++ {
		contentYaml := Root{}

		content, err := os.ReadFile("./mocks/test.yaml")

		if err != nil {
			b.Fatal(err)
		}

		Unmarshal(&contentYaml, string(content))

		if contentYaml.Config.Features.Caching.Host != "cache.example.com" {
			b.Errorf("BenchmarkUnmarshal() error to unmarshal file content")
		}
	}
}

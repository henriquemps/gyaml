# GYAML

`gyaml` is a Go library for parsing and unmarshaling YAML files into Go structs. It supports loading data from files or directly from strings.

## Features
- YAML parsing and reading
- Support for nested structures
- Automatic value conversion
- Array handling

## Installation

To install the library, use:

```sh
go get github.com/henriquemps/gyaml
```

## Usage

### Example 1: Reading YAML from a file

```go
package main

import (
    "fmt"
    "log"
    "github.com/henriquemps/gyaml"
)

type Config struct {
    Version string `yaml:"version"`
    App struct {
        Name string `yaml:"name"`
    } `yaml:"app"`
}

func main() {
    var config Config
    gyaml.FUnmarshal(&config, "config.yaml")
    fmt.Println("App Name:", config.App.Name)
}
```

### Example 2: Reading YAML from a raw string

```go
package main

import (
    "fmt"
    "log"
    "github.com/henriquemps/gyaml"
)

type Config struct {
    Version string `yaml:"version"`
    App struct {
        Name string `yaml:"name"`
    } `yaml:"app"`
}

func main() {
    yamlContent := `
    version: 1.0.0
    app:
      name: SampleApp
    `

    var config Config
    gyaml.Unmarshal(&config, yamlContent)
    fmt.Println("App Name:", config.App.Name)
}
```

## Main Functions

### `FUnmarshal(dt any, path string)`
Loads a YAML file and deserializes its content into a Go struct.

### `Unmarshal(dt any, v string)`
Reads and processes a YAML provided as a string.

## License

This project is licensed under the MIT license. See the LICENSE file for more details.

// Package gyaml provides functionality to parse and unmarshal YAML content into Go structs.
// It supports both loading YAML from a file or directly from a string input.
//
// The package offers two main functions for unmarshaling data:
//   - FUnmarshal: Unmarshals data from a YAML file specified by its path.
//   - Unmarshal: Reads and loads YAML content from a string and maps it to a Go struct.
//
// The package also includes various utility functions for handling YAML structures, such as
// extracting field values, handling arrays, and managing nested structures.
//
// Example usage:
//
//	import "path/to/gyaml"
//
//	var data MyStruct
//	err := gyaml.FUnmarshal(&data, "path/to/file.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Author: Henrique Matos
// Github: https://github.com/henriquemps/gyaml
// License: MIT
package gyaml

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Yaml struct {
	scanner *bufio.Scanner
	lines   []map[string]any
	path    string
	content string
}

var instance *Yaml

func newYaml() *Yaml {
	if instance == nil {
		instance = &Yaml{}
	}

	return instance
}

func (y *Yaml) init(path string, content string) {

	instance.path = path
	instance.content = content

	if y.path != "" {
		instance.scanner = getContentFile(path)
		instance.lines = getLinesFile(instance.scanner)
	} else {
		instance.scanner = bufio.NewScanner(strings.NewReader(content))
		instance.lines = getLinesFile(instance.scanner)
	}
}

// FUnmarshal allow Unmarshal per file
func FUnmarshal(dt any, path string) {
	i := newYaml()

	i.init(path, "")

	Unmarshal(dt, "")
}

// Unmarshal Read load file .yaml and build struct
func Unmarshal(dt any, v string) {

	if v != "" {
		i := newYaml()
		i.init("", v)
	}

	keys := make(map[string]any)
	var listValues []string

	lines := instance.lines

	for i := 0; i < len(lines); i++ {

		//keyPath := make([]string, 0)

		// Recover data of the current line
		currentText, currentField, currentValue, currentValueIsOnlyField, currentSpace := extractDataLine(lines, i)

		if currentValueIsOnlyField {
			continue
		}

		// Section to handle fields with items in the format:
		// campo:
		//   - A
		//   - B
		//   - C
		if isItemArray(currentText) {
			listValues = append(listValues, extractDataItemArray(currentText))

			// test if can we access next line
			if i+1 <= len(lines)-1 {

				// Recover data of the next line
				nextText, _, _, _, _ := extractDataLine(lines, i+1)

				if isItemArray(nextText) {
					continue
				}
			}
		}

		fullKey := buildKeyHierarchy(lines, currentSpace, i)

		if len(listValues) > 0 {
			keys[fullKey] = listValues
			listValues = nil
		} else {
			keys[fmt.Sprintf("%s.%s", fullKey, currentField)] = currentValue
		}
	}

	rootKeys := extractRootKeys(keys)
	structName := strings.ToLower(reflect.TypeOf(dt).Elem().Name())

	for key, value := range keys {
		setValueOnDataStruct(dt, key, value, 0, containsRootkey(rootKeys, structName))
	}
}

// buildKeyHierarchy responsible for building key hierarchy and values
//
// Example:
//
// config:
//
//	version: 1.0.0
//	app:
//	  name: SampleApp
//
// Hierarchy result for "name: SampleApp" is: config.app.name = SampleAPP
func buildKeyHierarchy(lines []map[string]any, currentSpace int, currentIndex int) string {

	keyPath := make([]string, 0)
	//keys := make(map[string]any)
	spaceLastField := lines[currentIndex-1]["spaces"].(int)

	for a := currentIndex - 1; a >= 0; a-- {

		// Recover data of the before line
		_, fieldPrevious, _, valuePreviousIsOnlyField, spacePreviousField := extractDataLine(lines, a)

		// if is first iteration then we can't check the spaces
		// in conditional
		if currentIndex-1 == a {
			if valuePreviousIsOnlyField && spacePreviousField < currentSpace {
				spaceLastField = spacePreviousField
				keyPath = prepend(keyPath, fieldPrevious)
			}
		} else {
			if valuePreviousIsOnlyField && spaceLastField > spacePreviousField && spacePreviousField < currentSpace {
				spaceLastField = spacePreviousField
				keyPath = prepend(keyPath, fieldPrevious)
			}
		}
	}

	fullKey := strings.Join(keyPath, ".")

	//if len(listValues) > 0 {
	//	keys[fullKey] = listValues
	//	listValues = nil
	//} else {
	//	fullKey = fmt.Sprintf("%s.%s", fullKey, currentField)
	//}

	return fullKey //, currentValue, listValues
}

// prepend add item to the beginning of the list
func prepend(list []string, item string) []string {
	return append([]string{item}, list...)
}

// containsRootkey check if a key exist in list
func containsRootkey(list []string, key string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

// extractRootKeys extract all root keys
func extractRootKeys(keys map[string]any) []string {
	rootKeys := make([]string, 0)

	for key, _ := range keys {
		fields := strings.Split(key, ".")

		if containsRootkey(rootKeys, fields[0]) {
			continue
		}

		rootKeys = append(rootKeys, fields[0])
	}

	return rootKeys
}

// extractDataLine extract data as field and values within then
func extractDataLine(lines []map[string]any, index int) (string, string, string, bool, int) {

	lineText := lines[index]["value"].(string)
	space := lines[index]["spaces"].(int)
	isField := strings.HasSuffix(lineText, ":")
	field, value, _ := strings.Cut(lineText, ":")
	field = strings.ToLower(strings.TrimSpace(field))
	value = strings.TrimSpace(value)

	return lineText, field, value, isField, space
}

// extractDataItemArray extract item value as array
func extractDataItemArray(value string) string {
	return strings.TrimSpace(strings.Replace(value, "-", "", -1))
}

// isItemArray check if the value is an item in an array
func isItemArray(value string) bool {
	compile, _ := regexp.Compile("-\\s*[a-zA-Z_\\-]+")

	return compile.MatchString(strings.TrimSpace(value))
}

// getContentFile open and get file content
func getContentFile(path string) *bufio.Scanner {

	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(file)

	var buffer bytes.Buffer

	newScanner := bufio.NewScanner(file)

	for newScanner.Scan() {
		buffer.WriteString(newScanner.Text() + "\n")
	}

	return bufio.NewScanner(&buffer)
}

// getLinesFile organize each line of the text in a structure
// each line is structured with data of the "space" and "value" of the line
func getLinesFile(scanner *bufio.Scanner) []map[string]any {
	var lines []map[string]any

	for scanner.Scan() {
		value := scanner.Text()

		// Ignore spaces
		if value == "" {
			continue
		}

		// Ignore commented line
		if []rune(value)[0] != '#' {
			totalSpaces := len(value) - len(strings.TrimLeft(value, " "))
			lines = append(lines, map[string]any{"spaces": totalSpaces, "value": scanner.Text()})
		}
	}

	return lines
}

// setValueOnDataStruct sets values within a struct.
// It navigates through the key hierarchy to find the corresponding field.
// When the field is found, it converts the value to the type allowed by the struct field and sets the value.
//
// In cases where there are multiple root structures in the YAML file, we need to determine
// whether we should remove the first "key" mapped in the hierarchy or not.
func setValueOnDataStruct(s any, key string, value any, indexField int, removeFirstKey bool) {

	structReflect := reflect.ValueOf(s).Elem()

	fields := strings.Split(key, ".")

	if removeFirstKey {
		removeFirstKey = true
		fields = fields[1:]
	}

	isLastField := indexField == len(fields)-1

	for i := 0; i < structReflect.NumField(); i++ {
		objectField := structReflect.Type().Field(i)
		objectValue := structReflect.Field(i)

		field := fields[indexField]

		if !isLastField && ((objectField.Tag.Get("yaml") == field) || (strings.ToLower(objectField.Name) == field)) {
			indexField++
			setValueOnDataStruct(objectValue.Addr().Interface(), key, value, indexField, removeFirstKey)
		} else if (isLastField && objectField.Tag.Get("yaml") == field) || (strings.ToLower(objectField.Name) == field) {
			objectValue.Set(reflect.ValueOf(convert(objectField.Type, value)))
			break
		}
	}
}

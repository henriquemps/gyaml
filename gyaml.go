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
	"unicode/utf8"
)

type Yaml struct {
	scanner   *bufio.Scanner
	lines     []map[string]any
	path      string
	content   string
	structure any
}

var instance *Yaml

func newYaml() *Yaml {
	if instance == nil {
		instance = &Yaml{}
	}

	return instance
}

func (y *Yaml) init(structure any, path string, content string) {

	instance.path = path
	instance.content = content
	instance.structure = structure

	if y.path != "" {
		instance.scanner = getContentFile(path)
		instance.lines = getLinesFile(instance.scanner)
	} else {
		instance.scanner = bufio.NewScanner(strings.NewReader(content))
		instance.lines = getLinesFile(instance.scanner)
	}
}

// FUnmarshal allow unmarshal per file
func FUnmarshal(structure any, path string) {

	i := newYaml()
	i.init(structure, path, "")

	build()
}

// Unmarshal allow unmarshal per raw string
func Unmarshal(structure any, rawContent string) {

	i := newYaml()
	i.init(structure, "", rawContent)

	build()
}

// build is the core logic for parsing and unmarshaling the YAML.
// It orchestrates the process of reading the lines from the YAML file or string,
// extracts data from each line, and organizes the information into a hierarchical
// key-value structure. The function also handles specific cases like arrays and
// multiline fields, updating the destination structure (passed as a parameter)
// with the values extracted from the YAML.
//
// The function is responsible for interpreting the structure of the YAML and
// ensuring that values are correctly mapped to the appropriate fields in the
// Go structure, respecting hierarchy and expected data types.
func build() {

	var listValues []string
	var strBuilder strings.Builder

	keys := make(map[string]any)
	enableMultiline := false

	for i := 0; i < len(instance.lines); i++ {

		// Recover data of the current line
		currentText, currentField, currentValue, currentValueIsOnlyField, _, currentSpace := extractDataLine(instance.lines, i)

		if currentValueIsOnlyField {
			continue
		}

		if ok := processItemArray(currentText, i, &listValues); ok {
			continue
		}

		if ok := processMultiline(currentText, i, &strBuilder, &enableMultiline); ok {
			continue
		}

		fullKey := buildKeyHierarchy(instance.lines, currentSpace, i)

		valueIsList := len(listValues) > 0
		valueIsTextMultiline := len(strBuilder.String()) > 0

		if valueIsList {
			keys[fullKey] = listValues
			listValues = nil
		}

		if valueIsTextMultiline {
			keys[fullKey] = strBuilder.String()
			strBuilder.Reset()
			enableMultiline = false
		}

		if !valueIsList && !valueIsTextMultiline {
			keys[fmt.Sprintf("%s.%s", fullKey, currentField)] = currentValue
		}
	}

	rootKeys := extractRootKeys(keys)
	structName := strings.ToLower(reflect.TypeOf(instance.structure).Elem().Name())

	for key, value := range keys {
		setValueOnDataStruct(instance.structure, key, value, 0, containsRootkey(rootKeys, structName))
	}
}

// processItemArray manager fields with items in the format:
//
// campo:
//   - A
//   - B
//   - C
func processItemArray(currentText string, index int, listValues *[]string) bool {

	ok := false

	if isItemArray(currentText) {
		*listValues = append(*listValues, strings.TrimSpace(strings.Replace(currentText, "-", "", -1)))

		onDataNextLine(index, func(nextLineText, nextField, nextValue string, nextIsCommonField, nextIsField bool, nextSpace int) {

			if isItemArray(nextLineText) {
				ok = true
			}
		})
	}

	return ok
}

// processMultiline manager fields multilines with operator ">" or "|", example:
//
//	 multiA: >
//	  Texto 1
//	  Texto 2
//	  Texto 3
//	  Texto 4
//
//	multiB: |
//	  Texto 1
//	  Texto 2
//	  Texto 3
//	  Texto 4
//
//	multiB: |-
//	  Texto 1
//	  Texto 2
//	  Texto 3
//	  Texto 4
func processMultiline(currentText string, index int, str *strings.Builder, enableMultiline *bool) bool {

	ok := false

	if isMultiLine(currentText) {
		*enableMultiline = true
		return true
	}

	if *enableMultiline {
		str.WriteString(currentText + "\n")

		onDataNextLine(index, func(nextLineText, nextField, nextValue string, nextIsCommonField, nextIsField bool, nextSpace int) {

			if !nextIsField {
				ok = true
			}
		})
	}

	return ok
}

// onDataNextLine checks if it can access the next line of the content
// and extracts the data from the next line
func onDataNextLine(index int, f func(string, string, string, bool, bool, int)) {

	if index+1 <= len(instance.lines)-1 {
		f(extractDataLine(instance.lines, index+1))
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
	spaceLastField := lines[currentIndex-1]["spaces"].(int)

	for a := currentIndex - 1; a >= 0; a-- {

		// Recover data of the before line
		_, fieldPrevious, _, _, isFieldPrevious, spacePreviousField := extractDataLine(lines, a)

		// if is first iteration then we can't check the spaces
		// in conditional
		if currentIndex-1 == a {
			if isFieldPrevious && spacePreviousField < currentSpace {
				spaceLastField = spacePreviousField
				keyPath = prepend(keyPath, fieldPrevious)
			}
		} else {
			if isFieldPrevious && spaceLastField > spacePreviousField && spacePreviousField < currentSpace {
				spaceLastField = spacePreviousField
				keyPath = prepend(keyPath, fieldPrevious)
			}
		}
	}

	fullKey := strings.Join(keyPath, ".")

	return fullKey
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
func extractDataLine(lines []map[string]any, index int) (string, string, string, bool, bool, int) {

	lineText := lines[index]["value"].(string)
	space := lines[index]["spaces"].(int)
	field, value, _ := strings.Cut(lineText, ":")
	isCommonField := strings.HasSuffix(lineText, ":")
	isField := isAnyField(lineText)
	field = strings.ToLower(strings.TrimSpace(field))
	value = strings.TrimSpace(value)

	return lineText, field, value, isCommonField, isField, space
}

// isItemArray check if the value is an item in an array
func isItemArray(value string) bool {
	compile, _ := regexp.Compile("-\\s*[a-zA-Z_\\-]+")

	return compile.MatchString(strings.TrimSpace(value))
}

// isMultiLine check if the value is multi line
func isMultiLine(value string) bool {
	compile, _ := regexp.Compile(".*:\\s*(\\|-|>|\\|)")

	return compile.MatchString(strings.TrimSpace(value))
}

// isMultiLine check if the value is multi line
func isAnyField(value string) bool {
	compile, _ := regexp.Compile(".*:\\s?(>|(\\|\\-?)?)")

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
			value = strings.Split(value, "#")[0]

			totalSpaces := utf8.RuneCountInString(value) - utf8.RuneCountInString(strings.TrimLeft(value, " "))
			lines = append(lines, map[string]any{"spaces": totalSpaces, "value": value})
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

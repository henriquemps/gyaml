package yaml

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// Read load file .yaml and build struct
func Read(dataStruct any, path string) error {

	scanner, _ := getContentFile(path)
	lines := getLinesFile(scanner)

	keys := make(map[string]any)
	listValues := make([]string, 0)

	for i := 0; i < len(lines); i++ {

		keyPath := make([]string, 0)

		// Recover data of the current line
		currentSpace := lines[i]["spaces"].(int)
		currentText, currentField, currentValue, currentValueIsOnlyField := extractDataLine(lines[i]["value"].(string))

		if !currentValueIsOnlyField {
			spaceLastField := lines[i-1]["spaces"].(int)

			if isItemArray(currentText) {

				listValues = append(listValues, extractDataItemArray(currentText))

				if i+1 <= len(lines)-1 {

					// Recover data of the next line
					nextText, _, _, _ := extractDataLine(lines[i+1]["value"].(string))

					if isItemArray(nextText) {
						continue
					}
				}
			}

			for a := i - 1; a >= 0; a-- {

				// Recover data of the before line
				spacePreviousField := lines[a]["spaces"].(int)
				_, fieldPrevious, _, valuePreviousIsOnlyField := extractDataLine(lines[a]["value"].(string))

				// if is first iteration then we can't check the spaces
				// in conditional
				if i-1 == a {
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

			if len(listValues) > 0 {
				keys[strings.Join(keyPath, ".")] = listValues
				listValues = make([]string, 0)
			} else {
				keys[fmt.Sprintf("%s.%s", strings.Join(keyPath, "."), currentField)] = currentValue
			}
		}
	}

	for key, value := range keys {
		setValueOnDataStruct(dataStruct, key, value, 0)
	}

	fmt.Println(dataStruct)

	return nil
}

// prepend add item to the beginning of the list
func prepend(list []string, item string) []string {
	return append([]string{item}, list...)
}

// extractDataLine extract data as field and values within then
func extractDataLine(line string) (string, string, string, bool) {

	lineText := line
	isField := strings.HasSuffix(lineText, ":")
	field, value, _ := strings.Cut(lineText, ":")
	field = strings.TrimSpace(field)
	value = strings.TrimSpace(value)

	return lineText, field, value, isField
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
func getContentFile(path string) (*bufio.Scanner, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var buffer bytes.Buffer

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		buffer.WriteString(scanner.Text() + "\n")
	}

	return bufio.NewScanner(&buffer), nil
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

func setValueOnDataStruct(s any, key string, value any, indexField int) {

	fields := strings.Split(key, ".")
	fields = fields[1:]
	isLastField := indexField == len(fields)-1

	sReflect := reflect.ValueOf(s).Elem()

	for i := 0; i < sReflect.NumField(); i++ {
		objectField := sReflect.Type().Field(i)
		objectValue := sReflect.Field(i)

		field := fields[indexField]

		if !isLastField && ((objectField.Tag.Get("yaml") == field) || (strings.ToLower(objectField.Name) == field)) {
			indexField++
			setValueOnDataStruct(objectValue.Addr().Interface(), key, value, indexField)
		} else if (isLastField && objectField.Tag.Get("yaml") == field) || (strings.ToLower(objectField.Name) == field) {
			objectValue.Set(reflect.ValueOf(convert(objectField.Type, value)))
			break
		}
	}
}

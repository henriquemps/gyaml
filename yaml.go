package yaml

import (
	"bufio"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// Read load file .yaml and build struct
func Read(dataStruct any, path string) error {

	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	buildStruct(dataStruct, lines)

	return nil
}

func buildStruct(dataStruct any, lines []string) {

	valuesArray := make([]string, 0)

	for index, line := range lines {
		field, _, _ := strings.Cut(strings.TrimSpace(line), ":")
		isOnlyField := strings.HasSuffix(line, ":")

		if !isOnlyField {
			regexpFieldValue := regexp.MustCompile(".*:.*")
			regexpItemArray := regexp.MustCompile("(-(\\s|)[a-zA-Z]+)|(-(\\s|)[0-9]+)")
			regexpItemDate := regexp.MustCompile("([0-9]{4,}-[0-9]{2,}-[0-9]{2,})")
			regexpArray := regexp.MustCompile(`\[.*]`)

			if regexpFieldValue.MatchString(line) || regexpItemDate.MatchString(line) {
				_, value, _ := strings.Cut(line, ":")
				setValueOnDataStruct(dataStruct, field, value)
			}

			if regexpItemArray.MatchString(line) {
				nextLine := index + 1

				line = strings.TrimSpace(strings.Replace(line, "-", "", -1))
				valuesArray = append(valuesArray, line)
				totalLines := len(lines)
				totalValues := len(valuesArray)

				if (nextLine < totalLines && !regexpItemArray.MatchString(lines[nextLine])) || nextLine == totalLines {
					field = lines[totalLines-totalValues-1]
					field = strings.Replace(field, ":", "", -1)

					setValueOnDataStruct(dataStruct, field, valuesArray)
				}
			}

			if regexpArray.MatchString(line) {
				strReplaced := strings.Trim(line, "[]")
				value := strings.Split(strReplaced, ",")

				setValueOnDataStruct(dataStruct, field, value)
			}
		} else {
			valuesArray = []string{}
		}
	}
}

// Build data within struct
func setValueOnDataStruct(s any, field string, value any) any {

	sReflect := reflect.ValueOf(s)

	if sReflect.Kind() == reflect.Ptr {
		sReflect = sReflect.Elem()
	}

	for i := 0; i < sReflect.NumField(); i++ {
		objectField := sReflect.Type().Field(i)
		objectValue := sReflect.Field(i)

		if (objectField.Tag.Get("yaml") == field) || (strings.ToLower(objectField.Name) == field) {
			objectValue.Set(reflect.ValueOf(convert(objectField.Type, value)))
		}

		if objectValue.Kind() == reflect.Struct {
			setValueOnDataStruct(objectValue.Addr().Interface(), field, value)
		}
	}

	return s
}

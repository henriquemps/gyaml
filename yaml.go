package yaml

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Read load file .yaml and build struct
func Read(dataStruct any, path string) error {

	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if !strings.HasSuffix(scanner.Text(), ":") {

			values := strings.Split(scanner.Text(), ":")

			buildStruct(dataStruct, values[0], values[1])
		}
	}

	return nil
}

// Build data within struct
func buildStruct(s any, field string, value any) any {

	sReflect := reflect.ValueOf(s)

	if sReflect.Kind() == reflect.Ptr {
		sReflect = sReflect.Elem()
	}

	for i := 0; i < sReflect.NumField(); i++ {
		objectField := sReflect.Type().Field(i)
		objectValue := sReflect.Field(i)

		if (objectField.Tag.Get("yaml") == field) || (strings.ToLower(objectField.Name) == strings.TrimSpace(field)) {
			objectValue.Set(reflect.ValueOf(convert(objectField.Type, value)))
		}

		if objectValue.Kind() == reflect.Struct {
			buildStruct(objectValue.Addr().Interface(), field, value)
		}
	}

	return s
}

// Convert data to a kind specific
func convert(convertTo reflect.Type, value any) any {

	result := value

	switch convertTo.Kind() {
	case reflect.Int:
		result, _ = strconv.Atoi(result.(string))
	case reflect.Bool:
		result, _ = strconv.ParseBool(result.(string))
	default:
		result = result.(string)
	}

	return result
}

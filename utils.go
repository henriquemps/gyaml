package gyaml

import (
	"reflect"
	"strconv"
)

// Convert data to a kind specific
func convert(convertTo reflect.Type, value any) any {

	result := value

	switch convertTo.Kind() {
	case reflect.Int:
		result, _ = strconv.Atoi(result.(string))
	case reflect.Bool:
		result, _ = strconv.ParseBool(result.(string))
	case reflect.Slice:
		result = result.([]string)
	default:
		result = result.(string)
	}

	return result
}

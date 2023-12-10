package settings

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Settings struct {
	CalendarID string
}

func getSettings() *Settings {
	// TODO: get settings from database
	return &Settings{}
}

func getSetting(settingPath string) (interface{}, error) {
	settings := getSettings()

	pathParts := strings.Split(settingPath, ".")
	val := reflect.ValueOf(settings).Elem()

	for _, part := range pathParts {
		// If it's an array, handle indexing
		if idx := strings.Index(part, "["); idx != -1 {
			arrayIndex, err := strconv.Atoi(part[idx+1 : len(part)-1])
			if err != nil {
				return nil, fmt.Errorf("invalid array index in %s", part)
			}
			part = part[:idx]

			val = val.FieldByName(part)
			if !val.IsValid() {
				return nil, fmt.Errorf("setting %s does not exist", part)
			}
			if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
				return nil, fmt.Errorf("%s is not an array", part)
			}
			if arrayIndex >= val.Len() {
				return nil, fmt.Errorf("index out of range in %s", part)
			}
			val = val.Index(arrayIndex)
		} else {
			val = val.FieldByName(part)
			if !val.IsValid() {
				return nil, fmt.Errorf("setting %s does not exist", part)
			}
		}
	}

	switch val.Kind() {
	case reflect.String:
		return val.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int(), nil
	case reflect.Bool:
		return val.Bool(), nil
	default:
		return nil, errors.New("unsupported type")
	}
}

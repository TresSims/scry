package facts

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var ErrUnfulfilledInterface = errors.New("unable to completely fulfill interface")

// ReadProcFileIntoStruct takes a path to a system file and read it into a struct
//
// The function uses `proc` tags to determine what key should fill each value of the
// struct and populate the struct. It returns when the struct is fill, and [ErrUnfulfilledInterface ]
// if the struct cannot be fully populated
func ReadProcFileIntoStruct[T any](path string, t *T) error {
	filePath, err := os.Open(path)
	if err != nil {
		return err
	}

	reflection := reflect.ValueOf(t)
	if reflection.Kind() == reflect.Pointer {
		reflection = reflection.Elem()
	}

	procToField := map[string]string{}
	for f := range reflection.Fields() {
		procToField[f.Tag.Get("proc")] = f.Name
	}

	fileScanner := bufio.NewScanner(filePath)

	for fileScanner.Scan() {

		k, v := getProcKV(fileScanner.Text())
		if fieldName, ok := procToField[k]; ok {
			field := reflection.FieldByName(fieldName)
			switch field.Type().Kind() {
			case reflect.String:
				field.SetString(v)
			case reflect.Int:
				i, err := strconv.ParseInt(v, 10, field.Type().Bits())
				if err != nil {
					return err
				}

				field.SetInt(i)
			case reflect.Bool:
				i, err := strconv.ParseBool(v)
				if err != nil {
					return err
				}

				field.SetBool(i)
			case reflect.Float64, reflect.Float32:
				i, err := strconv.ParseFloat(v, field.Type().Bits())
				if err != nil {
					return err
				}

				field.SetFloat(i)
			}
		}

		if !isZero(t) {
			return nil
		}
	}

	return ErrUnfulfilledInterface
}

func isZero(x interface{}) bool {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsZero() {
			return true
		}
	}

	return false
}

func getProcKV(line string) (string, string) {
	parts := strings.Split(line, ":")

	lineKey := strings.TrimSpace(parts[key])
	lineVal := strings.TrimSpace(parts[val])

	return lineKey, lineVal
}

package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(obj any) string {
	v := reflect.ValueOf(obj)
	sb := new(strings.Builder)
	for i := 0; i < v.NumField(); i++ {
		tag, ok := v.Type().Field(i).Tag.Lookup("properties")
		if !ok {
			continue
		}
		tagTokens := strings.Split(tag, ",")

		property := v.Field(i)

		isOmitEmpty := len(tagTokens) == 2 && tagTokens[1] == "omitempty"

		switch property.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value := strconv.Itoa(int(property.Int()))
			if isOmitEmpty && value == "0" {
				continue
			}
			sb.WriteString(tagTokens[0] + "=" + value + "\n")
		case reflect.String:
			if isOmitEmpty && property.String() == "" {
				continue
			}
			sb.WriteString(tagTokens[0] + "=" + escapePropertiesString(property.String()) + "\n")
		case reflect.Bool:
			if isOmitEmpty && !property.Bool() {
				continue
			}
			if property.Bool() {
				sb.WriteString(tagTokens[0] + "=true\n")
			} else {
				sb.WriteString(tagTokens[0] + "=false\n")
			}
		case reflect.Float32, reflect.Float64:
			if isOmitEmpty && property.Float() == 0 {
				continue
			}
			sb.WriteString(tagTokens[0] + "=" + fmt.Sprintf("%f\n", property.Float()))
		default:
			continue
		}
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

func escapePropertiesString(s string) string {
	if strings.HasPrefix(s, " ") {
		s = "\\" + s
	}
	return strings.NewReplacer("/", "//", "\n", "/\n").Replace(s)
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

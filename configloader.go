package configloader

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Load(result interface{}) error {
	return LoadFlag("config", result)
}

func LoadFlag(configLocationFlag string, result interface{}) error {
	configLocation := flag.String(configLocationFlag, "config.json", "configuration file")
	flag.Parse()
	return LoadLocation(*configLocation, result)
}

func LoadLocation(location string, result interface{}) (err error) {
	file, err := os.Open(location)
	if err != nil {
		log.Println("error on config load: ", err)
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(result)
	if err != nil {
		log.Println("invalid config json: ", err)
		return err
	}
	handleEnvironmentVars(result)
	return
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func fieldNameToEnvName(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToUpper(strings.Join(a, "_"))
}

func handleEnvironmentVars(config interface{}) {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	configType := configValue.Type()
	for index := 0; index < configType.NumField(); index++ {
		fieldName := configType.Field(index).Name
		envName := fieldNameToEnvName(fieldName)
		envValue := os.Getenv(envName)
		if envValue != "" {
			log.Println("use environment variable: ", envName, " = ", envValue)
			if configValue.FieldByName(fieldName).Kind() == reflect.Int64 {
				i, _ := strconv.ParseInt(envValue, 10, 64)
				configValue.FieldByName(fieldName).SetInt(i)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.String {
				configValue.FieldByName(fieldName).SetString(envValue)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Slice {
				val := []string{}
				for _, element := range strings.Split(envValue, ",") {
					val = append(val, strings.TrimSpace(element))
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(val))
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Map {
				value := map[string]string{}
				for _, element := range strings.Split(envValue, ",") {
					keyVal := strings.Split(element, ":")
					key := strings.TrimSpace(keyVal[0])
					val := strings.TrimSpace(keyVal[1])
					value[key] = val
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(value))
			}

		}
	}
}

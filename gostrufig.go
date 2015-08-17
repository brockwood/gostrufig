package gostrufig

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// CFGDEFAULT stands for the default value for the variable if nothing is provided in the config.
const CFGDEFAULT string = `cfg-def`

// CFGNAMESPACE if "true", will include this member variable's name as part of the namespace path in the configuration.
// For example,
//	type MyConfig struct {
//		Name    string `cfg-ns:"true"`
//		Number  int `cfg-ns:"true"`
//		Message string
// 	}
// example := MyConfig{"brockwood", 123, "This is an example."}
// When retrieving this config for an application called "testApp", it can be found at "/testApp/brockwood/123/"
const CFGNAMESPACE string = `cfg-ns`

var gsfdriver Driver

// Configuration response codes
const (
	CONFIGNOTFOUND int = 100
	CONFIGFOUND    int = 200
)

type Gostrufig struct {
	appName        string
	driverLocation string
	config         interface{}
}

// GetGoStruFig returns an instace of the configuration with the given name, driver and configuration URL.
func GetGostrufig(appname, location string) Gostrufig {
	return Gostrufig{appName: appname, driverLocation: location}
}

// RetrieveConfig applies the value from config to the input.
func (gsf *Gostrufig) RetrieveConfig(target interface{}) {
	gsf.config = target
	setInitialStructValues(gsf.config, gsf.appName)
	configpath := gsf.generateNameSpacePath()
	loadstatus := gsfdriver.Load(gsf.driverLocation, configpath)
	if loadstatus == CONFIGFOUND {
		setStructValues(gsf.config, gsf.appName, configpath, true)
	}
}

func (gsf *Gostrufig) generateNameSpacePath() string {
	var newNamespace bytes.Buffer
	newNamespace.WriteString("/" + gsf.appName)
	structData := reflect.ValueOf(gsf.config).Elem()
	structType := structData.Type()
	for fieldNum := 0; fieldNum < structData.NumField(); fieldNum++ {
		field := structType.Field(fieldNum)
		fieldData := structData.Field(fieldNum)
		structnamespace := field.Tag.Get(CFGNAMESPACE)
		if len(structnamespace) == 0 {
			continue
		}
		if fieldData.Type().String() != `string` {
			panic(fmt.Sprintf("Only strings can be used for namespace configuration. Fieldname %s of type %s is not supported.", field.Name, fieldData.Type()))
		}
		if len(fieldData.String()) == 0 {
			panic(fmt.Sprintf("Fieldname '%s' of the structure '%s' is tagged as part of this configuration's namespace but was zero length.",
				field.Name, reflect.TypeOf(gsf.config)))
		}
		newNamespace.WriteString("/" + fieldData.String())
	}
	return newNamespace.String()
}

func setInitialStructValues(target interface{}, envpreface string) {
	setStructValues(target, envpreface, "", false)
}

func setStructValues(target interface{}, envpreface, driverpreface string, searchDriver bool) {
	structData := reflect.ValueOf(target).Elem()
	structType := structData.Type()
	for fieldNum := 0; fieldNum < structData.NumField(); fieldNum++ {
		var decodeError error
		field := structType.Field(fieldNum)
		fieldData := structData.Field(fieldNum)
		possibleEnvName := strings.ToUpper(envpreface + "_" + field.Name)
		possibleEnvValue := os.Getenv(possibleEnvName)
		defaultValue := field.Tag.Get(CFGDEFAULT)
		var driverValue string
		var possibleDriverPath string
		if searchDriver {
			possibleDriverPath = driverpreface + "/" + field.Name
			driverValue = gsfdriver.Retrieve(possibleDriverPath)
		}
		if fieldData.Kind() == reflect.Struct {
			setStructValues(fieldData.Addr().Interface(), possibleEnvName, possibleDriverPath, searchDriver)
		} else if len(possibleEnvValue) > 0 {
			decodeError = setValue(&fieldData, possibleEnvValue)
		} else if len(driverValue) > 0 {
			decodeError = setValue(&fieldData, driverValue)
		} else if len(defaultValue) > 0 {
			decodeError = setValue(&fieldData, defaultValue)
		}
		if decodeError != nil {
			panic(fmt.Sprintf("Error parsing the default value %s for the field %s:  %s", defaultValue, field.Name, decodeError))
		}
	}
}

func setValue(targetValue *reflect.Value, valueString string) error {
	switch targetValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		targetint, err := strconv.ParseInt(valueString, 10, targetValue.Type().Bits())
		if err != nil {
			return err
		}
		targetValue.SetInt(targetint)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		targetuint, err := strconv.ParseUint(valueString, 10, targetValue.Type().Bits())
		if err != nil {
			return err
		}
		targetValue.SetUint(targetuint)
		break
	case reflect.Float32, reflect.Float64:
		targetfloat, err := strconv.ParseFloat(valueString, targetValue.Type().Bits())
		if err != nil {
			return err
		}
		targetValue.SetFloat(targetfloat)
		break
	case reflect.Bool:
		targetbool, err := strconv.ParseBool(valueString)
		if err != nil {
			return err
		}
		targetValue.SetBool(targetbool)
		break
	case reflect.String:
		targetValue.SetString(valueString)
		break
	}
	return nil
}

func RegisterDriver(driver Driver) {
	gsfdriver = driver
}

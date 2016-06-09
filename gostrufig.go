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

const DRIVERROOT string = `DRIVER_ROOT`

// Configuration response codes
const (
	CONFIGNOTFOUND int = 100
	CONFIGFOUND    int = 200
)

type Gostrufig struct {
	appName string
	driver  Driver
	config  interface{}
}

// GetGoStruFig returns an instace of the configuration with the given name, configuration location,
// and the driver you wish to use.  If your app just needs access to the environment and the tags
// in the struct, feel free to pass in <nil>.  Otherwise instantiate a Gostrufig driver and pass
// that in.
func GetGostrufig(appname, location string, driver Driver) Gostrufig {
	driverrootenv := strings.ToUpper(appname) + "_" + DRIVERROOT
	driverrootoverride := os.Getenv(driverrootenv)
	if len(driverrootoverride) > 0 {
		location = driverrootoverride
	}
	gsf := Gostrufig{appName: appname}
	if driver != nil {
		driver.SetRootPath(location)
		gsf.registerDriver(driver)
	}
	return gsf
}

// RetrieveConfig applies the value from config to the input.
func (gsf *Gostrufig) RetrieveConfig(target interface{}) {
	gsf.config = target
	gsf.setInitialStructValues(gsf.config, gsf.appName)
	configpath := gsf.generateNameSpacePath()
	loadstatus := CONFIGNOTFOUND
	if gsf.driver != nil {
		loadstatus = gsf.driver.Load(configpath)
	}
	if loadstatus == CONFIGFOUND {
		gsf.setStructValues(gsf.config, gsf.appName, configpath, true)
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

func (gsf *Gostrufig) setInitialStructValues(target interface{}, envpreface string) {
	gsf.setStructValues(target, envpreface, "", false)
}

func (gsf *Gostrufig) setStructValues(target interface{}, envpreface, driverpreface string, searchDriver bool) {
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
			driverValue = gsf.driver.Retrieve(possibleDriverPath)
		}
		if fieldData.Kind() == reflect.Struct {
			gsf.setStructValues(fieldData.Addr().Interface(), possibleEnvName, possibleDriverPath, searchDriver)
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
	case reflect.Slice:
		raw := valueString[1 : len(valueString)-1]
		parts := strings.Split(raw, " ")
		sliceType := targetValue.Type()
		entryType := sliceType.Elem()
		sliceValue := reflect.MakeSlice(sliceType, len(parts), len(parts))
		for ii, _ := range parts {
			entryValue := reflect.New(entryType).Elem()
			setValue(&entryValue, parts[ii])
			sliceValue.Index(ii).Set(entryValue)
		}
		targetValue.Set(sliceValue)
		break
	}
	return nil
}

func (gsf *Gostrufig) registerDriver(driver Driver) {
	gsf.driver = driver
}

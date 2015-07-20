package gostrufig

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNamespace(t *testing.T) {
	desiredNamespace := "/appname/developer/calculator/"
	type NameSpace struct {
		DecodeDir   string
		Environment string `cfg-ns:"true"`
		Timer       int
		Type        string `cfg-ns:"true"`
		TestTimeout float64
	}
	namespace := NameSpace{"~/decode", "developer", 300, "calculator", 12.54}
	newPath := generateNameSpacePath(&namespace, "appname")
	if newPath != desiredNamespace {
		t.Errorf("Expected namespace of %s, instead got namespace of %s\n", desiredNamespace, newPath)
	} else {
		t.Logf("Received namespace of %s\n", newPath)
	}
}

func TestStructDefaults(t *testing.T) {
	type InternalStruct struct {
		MySubInt  int64 `cfg-def:"9223372036854775807"`
		MySubBool bool  `cfg-def:"false"`
	}
	type EveryType struct {
		MyInt     int           `cfg-def:"-32"`
		MyInt8    int8          `cfg-def:"-128"`
		MyInt16   int16         `cfg-def:"-32768"`
		MyInt32   int32         `cfg-def:"-2147483648"`
		MyInt64   int64         `cfg-def:"-9223372036854775808"`
		MyUInt    uint          `cfg-def:"32"`
		MyUInt8   uint8         `cfg-def:"255"`
		MyUInt16  uint16        `cfg-def:"65535"`
		MyUInt32  uint32        `cfg-def:"4294967295"`
		MyUInt64  uint64        `cfg-def:"18446744073709551615"`
		MyBool    bool          `cfg-def:"true"`
		MyString  string        `cfg-def:"four score and seven years ago"`
		MyFloat32 float32       `cfg-def:"4.123456"`
		MyFloat64 float64       `cfg-def:"-4.123456789"`
		MyTime    time.Duration `cfg-def:"300"`
		MyStruct  InternalStruct
	}
	blankStruct := EveryType{}
	setInitialStructValues(&blankStruct, "c2fo")
	populatedStruct := EveryType{-32, -128, -32768, -2147483648, -9223372036854775808,
		32, 255, 65535, 4294967295, 18446744073709551615, true,
		`four score and seven years ago`, 4.123456, -4.123456789, 300,
		InternalStruct{9223372036854775807, false}}
	if reflect.DeepEqual(blankStruct, populatedStruct) {
		t.Log("Blank and populated structs are the same.")
	} else {
		t.Errorf("Comparison of blank and populated structs failed.")
	}
}

func TestStructEnv(t *testing.T) {
	type InternalEnvStruct struct {
		MySubInt  int64
		MySubBool bool
	}
	type EnvStruct struct {
		MyFloat32 float32
		MyFloat64 float64
		MyStruct  InternalEnvStruct
	}
	blankStruct := EnvStruct{}
	os.Setenv("C2FO_MYFLOAT32", "4.123456")
	os.Setenv("C2FO_MYFLOAT64", "-4.123456789")
	os.Setenv("C2FO_MYSTRUCT_MYSUBINT", "9223372036854775807")
	os.Setenv("C2FO_MYSTRUCT_MYSUBBOOL", "false")
	setInitialStructValues(&blankStruct, "c2fo")
	populatedStruct := EnvStruct{4.123456, -4.123456789,
		InternalEnvStruct{9223372036854775807, false}}
	if reflect.DeepEqual(blankStruct, populatedStruct) {
		t.Log("Blank and populated structs are the same.")
	} else {
		t.Errorf("Comparison of blank and populated structs failed.")
	}
}

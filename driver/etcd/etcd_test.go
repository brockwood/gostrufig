package etcd

import (
	"fmt"
	"github.com/brockwood/gostrufig"
	"reflect"
	"testing"
)

func TestEtcd(t *testing.T) {
	fmt.Println("Hopefully you have a populated Etcd server running on 2379...")

	type SubNameSpace struct {
		SubThingFloat float64
		SubThingBool  bool
	}

	type NameSpace struct {
		DecodeDir   string `cfg-def:"/home/user/decoder"`
		Environment string `cfg-ns:"true" cfg-def:"developer"`
		Timer       int
		Type        string
		TestTimeout float64
		SubInfo     SubNameSpace
	}
	ns := NameSpace{}

	populatedSubNs := SubNameSpace{
		SubThingFloat: 3.12345,
		SubThingBool:  true,
	}

	populatedNs := NameSpace{
		DecodeDir:   "/home/user/decoder",
		Environment: "developer",
		Timer:       3600,
		Type:        "calculator",
		TestTimeout: 3.14159276,
		SubInfo:     populatedSubNs,
	}
	gsfdriver := GetGostrufigDriver()
	gsf := gostrufig.GetGostrufig("appname", "http://localhost:2379", gsfdriver)
	gsf.RetrieveConfig(&ns)
	if reflect.DeepEqual(ns, populatedNs) {
		t.Log("Blank and populated structs are the same.")
	} else {
		t.Errorf("Comparison of blank and populated structs failed; is your Etcd server running?")
	}
}

package gostrufig

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEtcd(t *testing.T) {
	fmt.Println("Hopefully you have a populated Etcd server running on 2379...")
	type NameSpace struct {
		DecodeDir   string `cfg-def:"/home/user/decoder"`
		Environment string `cfg-ns:"true" cfg-def:"developer"`
		Timer       int
		Type        string
		TestTimeout float64
	}
	ns := NameSpace{}
	populatedNs := NameSpace{"/home/user/decoder", "developer", 3600, "calculator", 3.14159276}
	gsf := GetGoStruFig("c2fo", "etcd", "http://localhost:2379")
	gsf.RetrieveConfig(&ns)
	if reflect.DeepEqual(ns, populatedNs) {
		t.Log("Blank and populated structs are the same.")
	} else {
		t.Errorf("Comparison of blank and populated structs failed.")
	}
}

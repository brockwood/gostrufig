package gostrufig

import (
	"fmt"
	"reflect"

	"github.com/coreos/go-etcd/etcd"
)

type EtcdDriver struct {
	client   *etcd.Client
	response *etcd.Response
	path     string
}

func getEtcdDriver(machine string) Driver {
	newdriver := EtcdDriver{}
	machines := []string{machine}
	newdriver.client = etcd.NewClient(machines)
	return Driver(&newdriver)
}

// Load reads etcd configurations from the given path into EtcdDriver
func (ed *EtcdDriver) Load(configStorePath string) int {
	var err error
	ed.path = configStorePath
	ed.response, err = ed.client.Get(configStorePath, true, true)
	if err != nil {
		if etcderr, ok := err.(*etcd.EtcdError); ok {
			return etcderr.ErrorCode
		} else {
			panic(fmt.Sprintf("Error retrieving data from etcd:  %s", err.Error()))
		}
	}
	return CONFIGFOUND
}

// Populate copies data from the EtcdDriver into the given struct.
func (ed *EtcdDriver) Populate(targetStruct interface{}) {
	structData := reflect.ValueOf(targetStruct).Elem()
	structType := structData.Type()
	for fieldNum := 0; fieldNum < structData.NumField(); fieldNum++ {
		var decodeError error
		field := structType.Field(fieldNum)
		fieldData := structData.Field(fieldNum)
		etcdpath := ed.path + field.Name
		possibleValue := findEtcdNode(ed.response.Node.Nodes, etcdpath)
		decodeError = setValue(&fieldData, possibleValue)
		if decodeError != nil {
			panic(fmt.Sprintf("Error parsing the value %s return from Etcd for the field %s:  %s", possibleValue, field.Name, decodeError))
		}
	}
}

func findEtcdNode(nodes etcd.Nodes, name string) string {
	for _, node := range nodes {
		if node.Key == name {
			return node.Value
		}
	}
	return ""
}

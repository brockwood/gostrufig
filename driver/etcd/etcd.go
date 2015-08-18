package etcd

import (
	"fmt"
	"github.com/brockwood/gostrufig"
	"github.com/coreos/go-etcd/etcd"
)

type EtcdDriver struct {
	client   *etcd.Client
	response *etcd.Response
	path     string
}

func GetGostrufigDriver() gostrufig.Driver {
	newdriver := EtcdDriver{}
	return gostrufig.Driver(&newdriver)
}

func (ed *EtcdDriver) SetRootPath(rootpath string) {
	machines := []string{rootpath}
	ed.client = etcd.NewClient(machines)
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
	return gostrufig.CONFIGFOUND
}

// Populate copies data from the EtcdDriver into the given struct.
func (ed *EtcdDriver) Retrieve(path string) string {
	return findEtcdNode(ed.response.Node.Nodes, path)
}

func findEtcdNode(nodes etcd.Nodes, name string) string {
	var foundValue string
	for _, node := range nodes {
		if node.Dir {
			foundValue = findEtcdNode(node.Nodes, name)
		} else if node.Key == name {
			foundValue = node.Value
		}
		if len(foundValue) > 0 {
			break
		}
	}
	return foundValue
}

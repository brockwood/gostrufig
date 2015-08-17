package etcd

import (
	"fmt"
	"github.com/brockwood/gostrufig"
	"github.com/coreos/go-etcd/etcd"
)

func init() {
	newdriver := EtcdDriver{}
	gostrufig.RegisterDriver(gostrufig.Driver(&newdriver))
}

type EtcdDriver struct {
	client   *etcd.Client
	response *etcd.Response
	path     string
}

// Load reads etcd configurations from the given path into EtcdDriver
func (ed *EtcdDriver) Load(location, configStorePath string) int {
	if ed.client == nil {
		ed.loadEtcdClient()
	}
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

func (ed *EtcdDriver) loadEtcdClient() {
	machines := []string{"http://localhost:2379"}
	ed.client = etcd.NewClient(machines)
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

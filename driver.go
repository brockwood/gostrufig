package gostrufig

type Driver interface {
	SetRootPath(rootpath string)
	Load(configStorePath string) int
	Retrieve(name string) string
}

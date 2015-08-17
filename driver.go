package gostrufig

type Driver interface {
	Load(location, configStorePath string) int
	Retrieve(name string) string
}

package gostrufig

type Driver interface {
	Load(configStorePath string) int
	Populate(targetStruct interface{})
}

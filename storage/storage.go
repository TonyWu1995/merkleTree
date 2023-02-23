package storage

type Storage interface {
	Save(path string, v interface{}) error
	Load(path string, v interface{}) error
}

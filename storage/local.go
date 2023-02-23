package storage

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)

type LocalStorage struct {
	lock sync.Mutex
}

func Marshal(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
func Unmarshal(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (local *LocalStorage) Save(path string, v interface{}) error {
	local.lock.Lock()
	defer local.lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}
func (local *LocalStorage) Load(path string, v interface{}) error {
	local.lock.Lock()
	defer local.lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, v)
}

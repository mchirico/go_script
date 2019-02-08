package jsonconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// ReadJSON read and unmarshal
func ReadJSON(file string, v interface{}) error {

	data, err := ReadFile(file)
	if err != nil {
		log.Printf("Can't read file")
		return err
	}

	err = json.Unmarshal([]byte(data), &v)
	if err != nil {
		log.Printf("json.Unmarshal error")
		return err
	}
	return nil
}

// ReadFile generic read
func ReadFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return string(data), err
}

// WriteFile generic write
func WriteFile(data string, file string) (int, error) {
	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		return -1, err
	}

	n, err := f.WriteString(data)
	return n, err
}

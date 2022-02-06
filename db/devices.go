package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DevicesIndex struct {
	index    map[string]string
	revIndex map[string]string
}

func LoadDevices(path string) (*DevicesIndex, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)

	var result = make(map[string]string)
	json.Unmarshal(data, &result)

	var reversed = make(map[string]string, len(result))
	for k, v := range result {
		reversed[v] = k
	}

	return &DevicesIndex{result, reversed}, nil
}

func (d *DevicesIndex) Get(key string) (string, bool) {
	value, ok := d.revIndex[key]
	if !ok {
		value, ok = d.index[key]
	}
	return value, ok
}

package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DevicesJsonDB struct {
	mapping map[string]string
}

func LoadDevices(path string) (*DevicesJsonDB, error) {
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

	return &DevicesJsonDB{reversed}, nil
}

// func (d *DevicesJson) Get(key string) (string, bool) {
// 	value, ok := d.mapping[key]
// 	if !ok {
// 		value, ok = d.index[key]
// 	}
// 	return value, ok
// }

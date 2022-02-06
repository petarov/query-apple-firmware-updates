package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadDevices(path string) (map[string]string, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)

	var result = make(map[string]string)
	json.Unmarshal(data, &result)

	return result, nil
}

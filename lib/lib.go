package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/uuid"
)

// GetID --
func GetID() uuid.UUID {
	return uuid.New()
}

// CreateTimeStamp --
func CreateTimeStamp() string {
	return time.Now().Format("20060102150405")
}

type fn func(string, interface{}) interface{}

// ReadJSON --
func ReadJSON(dbmodel string, filepath string, f fn) ([]interface{}, interface{}, error) {
	jsonfile, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}

	var data []interface{}
	tableprops := f(dbmodel, nil)

	bytes, _ := ioutil.ReadAll(jsonfile)
	json.Unmarshal(bytes, &data)

	defer jsonfile.Close()

	return data, tableprops, nil
}

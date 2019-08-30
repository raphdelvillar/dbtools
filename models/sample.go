package models

import (
	"dbtools/lib"
	"encoding/json"

	uuid "github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

// SampleModel --
type SampleModel struct {
	ID        uuid.UUID     `sql:"id,type:uuid,pk"`
	IDBson    bson.ObjectId `json:"_id" bson:"_id" sql:"-"`
	IDHex     string        `sql:"id_bson"`
}

// GetSampleModel --
func GetSampleModel(data interface{}) interface{} {
	var sampleModel SampleModel
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &sampleModel)
	sampleModel.ID = lib.GetID()
	sampleModel.IDHex = bson.ObjectId(sampleModel.IDBson).Hex()

	return (*SampleModel)(&sampleModel)
}

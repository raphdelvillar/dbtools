package master

import (
	"dbtools/lib"
	"encoding/json"

	uuid "github.com/google/uuid"
)

// SampleType --
type SampleType struct {
	ID          uuid.UUID `sql:"id,type:uuid,pk"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}

// GetSampleType --
func GetSampleType(data interface{}) interface{} {
	var sampleType SampleType
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &sampleType)
	sampleType.ID = lib.GetID()

	return (*SampleType)(&sampleType)
}

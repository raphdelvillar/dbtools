package models

// GetStruct --
func GetStruct(name string, data interface{}) interface{} {
	switch name {

	case "samplemodel":
		return GetSampleModel(data)
	}

	return nil
}

package util

import (
	"bytes"
	"encoding/json"
)

func StructToBytesBuffer(givenStruct interface{}) (buffer bytes.Buffer, err error) {
	var stringInterfaceMap map[string]interface{}

	marshalledJson, err := json.Marshal(givenStruct)
	if err != nil {
		return buffer, err
	}

	//fmt.Println(string(marshalledJson))

	err = json.Unmarshal(marshalledJson, &stringInterfaceMap)
	if err != nil {
		return buffer, err
	}

	err = json.NewEncoder(&buffer).Encode(stringInterfaceMap)
	return buffer, err
}

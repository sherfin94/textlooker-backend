package elastic

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"textlooker-backend/util"
)

func (textQuery *TextQuery) Buffer() (bytesBuffer bytes.Buffer, err error) {
	return util.StructToBytesBuffer(textQuery)
}

func ParseResult(body io.ReadCloser) (queryResult QueryResult, err error) {
	err = json.NewDecoder(body).Decode(&queryResult)
	if err != nil {
		log.Fatal(err)
		return queryResult, err
	}
	return queryResult, err
}

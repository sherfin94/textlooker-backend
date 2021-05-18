package elastic

import (
	"encoding/json"
	"io"
	"log"
)

func ParseResult(body io.ReadCloser) (queryResult QueryResult, err error) {
	err = json.NewDecoder(body).Decode(&queryResult)
	if err != nil {
		log.Println(err)
		return queryResult, err
	}
	return queryResult, err
}

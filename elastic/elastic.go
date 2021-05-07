package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var Client *elasticsearch.Client

func Initiate() {
	var err error
	Client, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
}

func Save(index string, body interface{}) error {
	data, _ := json.Marshal(body)
	request := esapi.IndexRequest{
		Index:   index,
		Body:    strings.NewReader(string(data)),
		Refresh: "true",
	}

	response, err := request.Do(context.Background(), Client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer response.Body.Close()

	if response.IsError() {
		log.Printf("[%s] Error indexing document", response.Status())
		log.Println(response)
		return errors.New("could not index document")
	}

	return nil
}
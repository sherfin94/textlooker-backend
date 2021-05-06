package elastic

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var client *elasticsearch.Client

func Initiate() {
	var err error
	client, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
}

func Save(index string, body interface{}) {
	data, _ := json.Marshal(body)
	request := esapi.IndexRequest{
		Index:   index,
		Body:    strings.NewReader(string(data)),
		Refresh: "true",
	}
	response, err := request.Do(context.Background(), client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer response.Body.Close()

	if response.IsError() {
		log.Printf("[%s] Error indexing document", response.Status())
		log.Println(response)
	}
}

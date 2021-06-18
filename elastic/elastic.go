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

func Save(index string, body interface{}, ID string) (string, error) {
	data, _ := json.Marshal(body)
	request := esapi.IndexRequest{
		Index:      index,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
		DocumentID: ID,
	}

	response, err := request.Do(context.Background(), Client)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return string("error"), err
	}

	defer response.Body.Close()

	if response.IsError() {
		log.Printf("[%s] Error indexing document", response.Status())
		log.Println(response)
		return "", errors.New("could not index document")
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return "", err
	}
	return responseData["_id"].(string), nil
}

func Query(query TextQuery, index string) (queryResult QueryResult, err error) {

	queryBuffer, err := query.Buffer()

	if err != nil {
		log.Printf("Error encoding query: %s", err)
		return queryResult, err
	}
	var response *esapi.Response
	response, err = Client.Search(
		Client.Search.WithContext(context.Background()),
		Client.Search.WithIndex(index),
		Client.Search.WithBody(&queryBuffer),
	)

	if err != nil {
		log.Printf("Error getting response: %s", err)
		return queryResult, err
	}

	defer response.Body.Close()

	if response.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(response.Body).Decode(&e); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			return queryResult, errors.New("elasticsearch query failed")
		}
		return queryResult, err
	}

	if queryResult, err = ParseResult(response.Body); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return queryResult, err
	}

	return queryResult, err
}

package elastic

import (
	"context"
	"textlooker-backend/deployment"

	"github.com/olivere/elastic"
)

var Client *elastic.Client

func Initiate() {
	context := context.Background()

	Client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	_, _, err = Client.Ping(deployment.GetEnv("ELASTIC_URL")).Do(context)
	if err != nil {
		panic(err)
	}
}

func Save(index string, body interface{}, ID string) (result string, err error) {
	// data, _ := json.Marshal(body)
	// request := esapi.IndexRequest{
	// 	Index:      index,
	// 	Body:       strings.NewReader(string(data)),
	// 	Refresh:    "true",
	// 	DocumentID: ID,
	// }

	// response, err := request.Do(context.Background(), Client)
	// if err != nil {
	// 	log.Fatalf("Error getting response: %s", err)
	// }

	// defer response.Body.Close()

	// if response.IsError() {
	// 	log.Printf("[%s] Error indexing document", response.Status())
	// 	log.Println(response)
	// 	return "", errors.New("could not index document")
	// }

	// var responseData map[string]interface{}
	// if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
	// 	log.Printf("Error parsing the response body: %s", err)
	// 	return "", err
	// }
	// return responseData["_id"].(string), nil
	return result, err
}

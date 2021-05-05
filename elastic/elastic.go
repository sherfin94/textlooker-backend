package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var client *elasticsearch.Client

func Initiate() {
	var err error
	client, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
}

package nlp

import (
	"log"

	"github.com/jdkato/prose/v2"
)

type Entity struct {
	Type string
	Text string
}

func ExtractEntities(text string) (entities []Entity) {
	if doc, err := prose.NewDocument(text); err != nil {
		log.Fatal(err)
	} else {
		for _, ent := range doc.Entities() {
			entities = append(entities, Entity{Type: ent.Label, Text: ent.Text})
		}
	}
	return entities
}

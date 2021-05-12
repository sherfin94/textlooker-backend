package nlp

import (
	"log"

	"github.com/bbalet/stopwords"
	"github.com/jdkato/prose/v2"
)

func Tokenize(text string) (tokens []string, err error) {
	cleanText := stopwords.CleanString(text, "en", true)
	doc, err := prose.NewDocument(cleanText)
	if err != nil {
		log.Fatal(err)
		return tokens, err
	}
	tokens = []string{}

	for _, token := range doc.Tokens() {
		if token.Text != "i" {
			tokens = append(tokens, token.Text)
		}
	}

	return tokens, nil
}

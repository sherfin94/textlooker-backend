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

func Ngrams(text string, n int) (ngrams []string, err error) {
	tokens, err := Tokenize(text)
	if err != nil {
		log.Fatal(err)
		return ngrams, err
	}

	for i := 0; i < len(tokens)-n+1; i++ {
		ngram := tokens[i]
		for j := 1; j < n; j++ {
			ngram += "." + tokens[i+j]
		}
		ngrams = append(ngrams, ngram)
	}

	return ngrams, err
}

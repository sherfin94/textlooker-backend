package nlp

import (
	"errors"
	"log"

	"github.com/bbalet/stopwords"
	"github.com/jdkato/prose/v2"
)

func getWords(text string) (tokens []string, err error) {
	cleanText := stopwords.CleanString(text, "en", true)
	doc, err := prose.NewDocument(cleanText)
	if err != nil {
		log.Println(err)
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

func ngrams(tokens []string, n int) (ngrams []string, err error) {

	for i := 0; i < len(tokens)-n+1; i++ {
		ngram := tokens[i]
		for j := 1; j < n; j++ {
			ngram += "." + tokens[i+j]
		}
		ngrams = append(ngrams, ngram)
	}

	return ngrams, err
}

func Tokenize(text string) (words []string, err error) {
	tokens, err := getWords(text)
	if err != nil {
		return tokens, err
	}

	unigrams, err1 := ngrams(tokens, 1)
	bigrams, err2 := ngrams(tokens, 2)
	trigrams, err3 := ngrams(tokens, 3)

	if err1 != nil || err2 != nil || err3 != nil {
		return tokens, errors.New("Unable to generate ngrams")
	}

	words = append(words, unigrams...)
	words = append(words, bigrams...)
	words = append(words, trigrams...)

	return words, nil
}

package nlp

import (
	sentimentAnalytics "github.com/cdipaolo/sentiment"
)

func Sentiment(text string) (sentiment string, err error) {
	model, err := sentimentAnalytics.Restore()
	if err != nil {
		return sentiment, err
	}

	analysis := model.SentimentAnalysis(text, sentimentAnalytics.English)

	if analysis.Score == 0 {
		sentiment = "negative"
	} else {
		sentiment = "positive"
	}

	return sentiment, err
}

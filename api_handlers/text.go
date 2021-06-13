package apihandlers

import (
	"errors"
	"textlooker-backend/models"
	"time"
)

func TextWithDate(content string, author []string, date time.Time, source *models.Source) (err error) {

	if !source.DateAvailable {
		return errors.New("this source does not have date enabled, please create a source with date enabled")
	}

	var text models.Text
	text, err = models.NewText(content, author, date, int(source.ID))

	if err != nil {
		return err
	} else {
		go text.SendToProcessQueue()
	}

	return err
}

func TextWithoutDate(content string, author []string, source *models.Source) (err error) {
	if source.DateAvailable {
		return errors.New("this source has date enabled, so date must be provided along with data")
	}

	var text models.Text
	text, err = models.NewTextWithoutDate(content, author, int(source.ID))

	if err != nil {
		return err
	} else {
		go text.SendToProcessQueue()
	}

	return err
}

package handlers

import (
	"errors"
	"strconv"
	"textlooker-backend/models"
	"textlooker-backend/util"
	"time"
)

func Text(content string, author []string, date string, source *models.Source) (err error) {
	var text models.Text
	var dateAsInteger int64
	if len(date) > 0 {
		if !source.DateAvailable {
			return errors.New("this source does not have date enabled, please create a source with date enabled")
		}
		dateAsInteger, err = strconv.ParseInt(date, 10, 64)
		time := util.ParseTimestamp(float64(dateAsInteger))
		if err != nil {
			return err
		}
		text, err = models.NewText(content, author, *time, int(source.ID))
	} else {
		now := time.Now()
		if source.DateAvailable {
			return errors.New("this source has date enabled, so date must be provided along with data")
		}
		text, err = models.NewText(content, author, now, int(source.ID))
	}

	if err != nil {
		return err
	} else {
		go text.SendToProcessQueue()
	}

	return err
}

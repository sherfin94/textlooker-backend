package handlers

import (
	"strconv"
	"textlooker-backend/models"
	"textlooker-backend/util"
)

func Text(content string, author []string, date string, sourceID int) (err error) {
	var text models.Text
	var dateAsInteger int64
	if len(date) > 0 {
		dateAsInteger, err = strconv.ParseInt(date, 10, 64)
		time := util.ParseTimestamp(float64(dateAsInteger))
		if err != nil {
			return err
		}
		text, err = models.NewText(content, author, *time, int(sourceID))
	} else {
		text, err = models.NewTextWithoutDate(content, author, int(sourceID))
	}

	if err != nil {
		return err
	} else {
		go text.SendToProcessQueue()
	}

	return err
}

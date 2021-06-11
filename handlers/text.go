package handlers

import (
	"strconv"
	"textlooker-backend/models"
	"textlooker-backend/util"
)

func Text(content string, author []string, date string, sourceID int) (err error) {
	dateAsInteger, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		return err
	}
	time := util.ParseTimestamp(float64(dateAsInteger))

	if text, err := models.NewText(content, author, *time, int(sourceID)); err != nil {
		return err
	} else {
		go text.SendToProcessQueue()
	}

	return err
}

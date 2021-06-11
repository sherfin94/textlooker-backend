package handlers

import (
	"errors"
	"strconv"
	"textlooker-backend/database"
	"textlooker-backend/models"
	"textlooker-backend/util"
)

func Text(content string, author []string, date string, sourceID int, user *models.User) (err error) {
	var source models.Source
	dateAsInteger, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		return err
	}
	time := util.ParseTimestamp(float64(dateAsInteger))
	database.Database.Where("user_id = ? and id = ?", user.ID, sourceID).Find(&source)
	if source.ID == 0 {
		err = errors.New("source not found")
		return err
	}

	if text, err := models.NewText(content, author, *time, int(source.ID)); err != nil {
		return err
	} else {
		go text.SendToProcessQueue()
	}

	return err
}

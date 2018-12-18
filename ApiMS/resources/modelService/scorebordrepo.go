package modelService

import (
	"github.com/go-park-mail-ru/2018_2_YetAnotherGame/ApiMS/resources/models"

	"github.com/jinzhu/gorm"
)

func GetScoreboardPage(db *gorm.DB, numberOfPage, countOfString int) (models.Scoreboard, bool) {
	var scoreboard models.Scoreboard
	db.Table("users").Order("score DESC").Offset(numberOfPage * countOfString).Limit(countOfString + 3).Find(&scoreboard.Users)
	if len(scoreboard.Users) > countOfString {
		return scoreboard, true
	}
	return scoreboard, false
}

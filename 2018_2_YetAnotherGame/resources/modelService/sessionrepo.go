package modelService

import (
	"2018_2_YetAnotherGame/resources/models"

	"github.com/jinzhu/gorm"
)

func GetSessionByEmail(db *gorm.DB, email string) models.Session {
	var tmp models.Session
	db.Table("sessions").Select("id, email").Where("email = ?", email).Scan(&tmp)
	return tmp
}

func GetSessionByID(db *gorm.DB, id string) models.Session {
	var tmp models.Session
	db.Table("sessions").Select("id, email").Where("id = ?", id).Scan(&tmp)
	return tmp
}

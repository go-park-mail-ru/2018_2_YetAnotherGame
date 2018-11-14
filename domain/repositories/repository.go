package repositories

import (
	"2018_2_YetAnotherGame/domain/models"

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

func GetScoreboardPage(db *gorm.DB, numberOfPage, countOfString int) ([]models.User, bool) {
	var users []models.User
	db.Table("users").Order("score DESC").Offset(numberOfPage * countOfString).Limit(countOfString + 3).Find(&users)
	if len(users) > countOfString {
		return users, true
	}
	return users, false
}

func FindUserByID(db *gorm.DB, id string) models.User {
	var user models.User
	db.Where("id = ?", id).First(&user)
	return user
}

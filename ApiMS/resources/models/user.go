package models

type User struct {
	ID        string `gorm:"primary_key"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Score     int    `json:"score"`
	Avatar    string `json:"avatar"`
	// Avatar    string `json:"photo_100"`
}

package models


type User struct {
	Email string `json:"email"`
	First_name string `json:"first_name"`
	Last_name string `json:"last_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Score int `json:"score"`
}
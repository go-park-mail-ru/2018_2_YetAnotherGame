package controllers

import (
	"2018_2_YetAnotherGame/domain/models"
	"2018_2_YetAnotherGame/domain/repositories"
	"2018_2_YetAnotherGame/domain/viewmodels"
	"2018_2_YetAnotherGame/infostructures/functions"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

func (env *Environment) RegistrationHandle(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	// userAvatar, handler, err := r.FormFile("image")

	session := repositories.GetSessionByEmail(env.DB, user.Email)

	if session.ID != "" {
		err := functions.BadRequest("already exists", w)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	ID := uuid.New()
	user.ID = ID.String()
	session.ID = ID.String()
	session.Email = user.Email
	env.DB.Create(&session)
	env.DB.Create(&user)

	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   ID.String(),
		Expires: time.Now().Add(60 * time.Minute),
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	message, err := json.Marshal(ID.String())
	if err != nil {
		logrus.Error(err)
	}
	w.Write(message)
}

func (env *Environment) LoginHandle(w http.ResponseWriter, r *http.Request) {
	authUser := models.Auth{}
	json.NewDecoder(r.Body).Decode(&authUser)
	if authUser.Password == "" || authUser.Email == "" {
		err := functions.BadRequest("Не указан E-Mail или пароль", w)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	session := repositories.GetSessionByEmail(env.DB, authUser.Email)
	if session.ID == "" {
		err := functions.BadRequest("Неверный E-Mail", w)
		if err != nil {
			logrus.Error(err)
		}
		return
	}
	user := repositories.FindUserByID(env.DB, session.ID)

	if user.Password != authUser.Password {
		err := functions.BadRequest("неверный пароль", w)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   session.ID,
		Expires: time.Now().Add(60 * time.Minute),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)

	message, err := json.Marshal(session.ID)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(message)
}

func (env *Environment) MeHandle(w http.ResponseWriter, r *http.Request) {
	Cookies, err := r.Cookie("sessionid")
	if err != nil {
		logrus.Warn("no cookies")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ID := Cookies.Value
	session := repositories.GetSessionByID(env.DB, ID)
	if session.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := repositories.FindUserByID(env.DB, session.ID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	message, err := json.Marshal(user)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(message)
}

func (env *Environment) ScoreboardHandle(w http.ResponseWriter, r *http.Request) {
	numberOfPage := 0
	countOfString := 3

	if r.URL.Query().Get("page") != "" {
		numberOfPage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	}

	scoreboard, canNext := repositories.GetScoreboardPage(env.DB, numberOfPage, countOfString)

	b := viewmodels.ScoreboardPageViewModel{}
	b.Scoreboard.Users = scoreboard.Users[:countOfString]
	b.CanNext = canNext
	message, err := json.Marshal(b)
	if err != nil {
		logrus.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(message)
}

func (env *Environment) LogOutHandle(w http.ResponseWriter, r *http.Request) {
	sessionCookies, _ := r.Cookie("sessionid")
	id := sessionCookies.Value
	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   id,
		Expires: time.Now(),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	message, err := json.Marshal(id)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(message)
}

func (env *Environment) AvatarHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		logrus.Error(err)
		return
	}
	defer file.Close()

	id, err := r.Cookie("sessionid")
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Println(id.Value)

}

func saveFile(db *gorm.DB, w http.ResponseWriter, file multipart.File, user_id string, handle *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Error(err)
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// fmt.Println(pwd)
	// src := pwd + "/uploads/" + users[user_id].Email + ".jpeg"
	tmpuser := models.User{}
	//db.Table("users ").Select("id, email, first_name, last_name, username ").Where("id = ?", tmpuser).Scan(&tmpuser)
	src := pwd + "/uploads/" + tmpuser.Email + handle.Filename

	err = ioutil.WriteFile(src, data, 0666)
	if err != nil {
		logrus.Error(err)
		return
	}
	//tmpuser.Avatar = src
	db.Table("users ").Select("id, email, first_name, last_name, username, password, score, avatar ").Where("id = ?", user_id).Scan(&tmpuser)
	//users.Store(user_id, tmpuser)
	db.Model(&tmpuser).Updates(models.User{Avatar: src}).Where("id = ?", user_id)
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

func (env *Environment) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	cookies, _ := r.Cookie("sessionid")
	id := cookies.Value
	userID := id
	//fmt.Println(r.MatchString("peach"))
	session := models.Session{}
	env.DB.Table("sessions").Select("id, email").Where("email = ?", userID).Scan(&session)

	if session.ID == "" {
		err := functions.BadRequest("Нет пользователей", w)
		if err != nil {
			logrus.Error(err)
		}
	}
	tmpuser := models.User{}
	env.DB.Table("users ").Select("id, email, first_name, last_name, username, password, score, avatar ").Where("id = ?", userID).Scan(&tmpuser)
	env.DB.Model(&tmpuser).Updates(models.User{Email: user.Email, FirstName: user.FirstName, LastName: user.LastName, Username: user.Username}).Where("id = ?", userID)
	env.DB.Save(&tmpuser)

}

/*	Берем данные которые приходят JSON'ом в теле ответа и парсим их в
 *	VKRosonseData добавляем к этому email и все эти данные пушим в проверяя нет ли такого юзера
 *	!!!Проблема: в вк может и не быть email'a, и тогда поле email в базе будет пустым, а оно ключевое
 *	Вариант решения: просить пользователя самому его вводить
 *
 */

const (
	APP_ID          = "6752650"
	APP_KEY         = "GUYoUbMgTZpYopPzrO5b"
	APP_SECRET      = "035ac1d8035ac1d8035ac1d8d4033dc8520035a035ac1d858b7a9b5f658c1e4bdba9b12"
	API_URL         = "https://api.vk.com/method/users.get?fields=email,photo_50&access_token=%s&v=5.52"
	API_RedirectURL = "http://127.0.0.1:8000/api/vkauth"
)

type VKResponseData struct {
	Response []struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Photo     string `json:"photo_100"`
	}
	Email string
}

// https://oauth.vk.com/authorize?client_id=6752650&redirect_uri=http://127.0.0.1:8000/api/vkauth&scope=4194306

func (env *Environment) VKRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.FormValue("code")
	conf := oauth2.Config{
		ClientID:     APP_ID,
		ClientSecret: APP_KEY,
		RedirectURL:  API_RedirectURL,
		Endpoint:     vk.Endpoint,
	}

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Printf("Bad Exchange: %v", err)
		return
	}
	var email string
	if token.Extra("email") != nil {
		email = token.Extra("email").(string)
	}

	client := conf.Client(ctx, token)

	resp, err := client.Get(fmt.Sprintf(API_URL, token.AccessToken))
	if err != nil {
		log.Printf("Bad resp: %v", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	if err != nil {
		log.Printf("Bad read body: %v", err)
		return
	}

	data := VKResponseData{}
	json.Unmarshal(body, &data)
	data.Email = email //Теперь из Имени и Фамилии можно сделать Юзернейм и не будет только пароля, но можно добавить токен для валидации

}

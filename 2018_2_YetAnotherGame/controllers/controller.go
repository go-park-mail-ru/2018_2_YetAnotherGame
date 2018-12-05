package controllers

import (
	"2018_2_YetAnotherGame/resources/models"
	"2018_2_YetAnotherGame/resources/modelService"

	"2018_2_YetAnotherGame/infostructures/functions"
	"fmt"
	"github.com/mailru/easyjson"
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

func (env *Environment) RegistrationHandle(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	//json.NewDecoder(r.Body).Decode(&user)
	easyjson.UnmarshalFromReader(r.Body, &user)
	// userAvatar, handler, err := r.FormFile("image")

	session := modelService.GetSessionByEmail(env.DB, user.Email)

	if session.ID != "" {
		err := functions.SendStatus("already exists", w, 400)
		env.Counter.WithLabelValues(r.URL.Path, "400").Inc()
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
	err := functions.SendStatus(ID.String(), w, 201)
	env.Counter.WithLabelValues(r.URL.Path, "201").Inc()
	if err != nil {
		logrus.Error(err)
	}
}

func (env *Environment) LoginHandle(w http.ResponseWriter, r *http.Request) {
	authUser := models.Auth{}
	//json.NewDecoder(r.Body).Decode(&authUser)
	easyjson.UnmarshalFromReader(r.Body, &authUser)
	if authUser.Password == "" || authUser.Email == "" {
		err := functions.SendStatus("Не указан E-Mail или пароль", w, 400)
		env.Counter.WithLabelValues(r.URL.Path, "400").Inc()
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	session := modelService.GetSessionByEmail(env.DB, authUser.Email)
	if session.ID == "" {
		err := functions.SendStatus("Неверный E-Mail", w, 400)
		env.Counter.WithLabelValues(r.URL.Path, "400").Inc()
		if err != nil {
			logrus.Error(err)
		}
		return
	}
	user := modelService.FindUserByID(env.DB, session.ID)

	if user.Password != authUser.Password {
		err := functions.SendStatus("неверный пароль", w, 400)
		env.Counter.WithLabelValues(r.URL.Path, "400").Inc()
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
	err := functions.SendStatus(session.ID, w, 201)
	if err != nil {
		logrus.Error(err)
	}
}

func (env *Environment) MeHandle(w http.ResponseWriter, r *http.Request) {
	Cookies, err := r.Cookie("sessionid")
	if err != nil {
		logrus.Warn("no cookies")
		w.WriteHeader(http.StatusUnauthorized)
		env.Counter.WithLabelValues(r.URL.Path, "401").Inc()
		return
	}
	ID := Cookies.Value
	session := modelService.GetSessionByID(env.DB, ID)
	if session.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		env.Counter.WithLabelValues(r.URL.Path, "400").Inc()
		return
	}
	user := modelService.FindUserByID(env.DB, session.ID)
	user.Password = "жулик, не воруй"
	w.WriteHeader(http.StatusCreated)
	env.Counter.WithLabelValues(r.URL.Path, "201").Inc()
	w.Header().Set("Content-Type", "application/json")
	message, err := easyjson.Marshal(user)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(message)
}

func (env *Environment) ScoreboardHandle(w http.ResponseWriter, r *http.Request, ) {
	numberOfPage := 0
	countOfString := 3

	if r.URL.Query().Get("page") != "" {
		numberOfPage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	}

	scoreboard, canNext := modelService.GetScoreboardPage(env.DB, numberOfPage, countOfString)

	b := models.ScoreboardPageViewModel{}
	b.Scoreboard.Users = scoreboard.Users[:countOfString]
	b.CanNext = canNext
	message, err := easyjson.Marshal(b)
	if err != nil {
		logrus.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
	env.Counter.WithLabelValues(r.URL.Path, "200").Inc()
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
	err := functions.SendStatus(id, w, 201)
	env.Counter.WithLabelValues(r.URL.Path, "201").Inc()
	if err != nil {
		logrus.Error(err)
	}
}
func (env *Environment) AvatarHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	file, handle, err := r.FormFile("image")
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

	saveFile(env.DB, w, file, id.Value, handle)
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

	tmpuser := models.User{}
	src := pwd + "/uploads/" + tmpuser.Email + handle.Filename

	err = ioutil.WriteFile(src, data, 0666)
	if err != nil {
		logrus.Error(err)
		return
	}
	db.Table("users ").Select("id, email, first_name, last_name, username, password, score, avatar ").Where("id = ?", user_id).Scan(&tmpuser)
	db.Model(&tmpuser).Updates(models.User{Avatar: src}).Where("id = ?", user_id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	msg := models.Message{"File uploaded successfully!"}
	message, err := easyjson.Marshal(msg)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(message)
}

func (env *Environment) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	//json.NewDecoder(r.Body).Decode(&user)
	easyjson.UnmarshalFromReader(r.Body, &user)
	cookies, _ := r.Cookie("sessionid")
	id := cookies.Value
	userID := id
	session := models.Session{}
	env.DB.Table("sessions").Select("id, email").Where("email = ?", userID).Scan(&session)

	if session.ID == "" {
		err := functions.SendStatus("Нет пользователей", w, 400)
		env.Counter.WithLabelValues(r.URL.Path, "400").Inc()
		if err != nil {
			logrus.Error(err)
		}
	}
	tmpuser := models.User{}
	env.DB.Table("users ").Select("id, email, first_name, last_name, username, password, score, avatar ").Where("id = ?", userID).Scan(&tmpuser)
	env.DB.Model(&tmpuser).Updates(models.User{Email: user.Email, FirstName: user.FirstName, LastName: user.LastName, Username: user.Username}).Where("id = ?", userID)
	env.DB.Save(&tmpuser)
	err := functions.SendStatus("Successfull", w, 201)
	env.Counter.WithLabelValues(r.URL.Path, "201").Inc()
	if err != nil {
		logrus.Error(err)
	}
}

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
	err = easyjson.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	data.Email = email //Теперь из Имени и Фамилии можно сделать Юзернейм и не будет только пароля, но можно добавить токен для валидации
	//TODO
}

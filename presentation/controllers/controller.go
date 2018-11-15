package controllers

import (
	"2018_2_YetAnotherGame/domain/models"
	"2018_2_YetAnotherGame/domain/repositories"
	"2018_2_YetAnotherGame/domain/viewmodels"
	"2018_2_YetAnotherGame/infostructures/functions"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
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
	user_id := id.Value

	saveFile(db, w, file, user_id, handle)
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

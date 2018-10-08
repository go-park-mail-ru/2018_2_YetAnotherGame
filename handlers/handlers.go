package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"../models"
)

//SignUp ..
func SignUp(ids map[string]string, users map[string]*models.User, w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	if _, ok := ids[user.Email]; ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Error{Msg: "already exists"}
		message, _ := json.Marshal(msg)
		w.Write(message)
	}

	id, _ := exec.Command("uuidgen").Output()

	stringID := string(id[:])
	stringID = strings.Trim(stringID, "\n")
	users[stringID] = &user
	ids[user.Email] = stringID

	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   stringID,
		Expires: time.Now().Add(60 * time.Minute),

		Path: "/",
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	message, _ := json.Marshal(stringID)
	w.Write(message)
}

//Login ..
func Login(ids map[string]string, users map[string]*models.User, w http.ResponseWriter, r *http.Request) {
	cred := models.Auth{}
	json.NewDecoder(r.Body).Decode(&cred)
	user_id := ids[cred.Email]
	if cred.Password == "" || cred.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Error{Msg: "Не указан E-Mail или пароль"}
		message, _ := json.Marshal(msg)
		w.Write(message)
	}

	if _, ok := users[user_id]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal("Неверный E-Mail")
		w.Write(message)
	}
	if users[user_id].Password != cred.Password {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Error{Msg: "Неверный пароль"}
		message, _ := json.Marshal(msg)

		w.Write(message)
	}

	ids[cred.Email] = user_id
	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   user_id,
		Expires: time.Now().Add(60 * time.Minute),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)

	message, _ := json.Marshal(user_id)
	w.Write(message)

}

//Me ..
func Me(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {
	id, err := r.Cookie("sessionid")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id2 := id.Value

	if _, ok := users[id2]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users[id2].Score++

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	message, _ := json.Marshal(users[id2])
	w.Write(message)
}

//Update ..
func Update(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {
	tmp := models.User{}
	json.NewDecoder(r.Body).Decode(&tmp)
	id, _ := r.Cookie("sessionid")
	id2 := id.Value
	user_id := id2
	//fmt.Println(r.MatchString("peach"))

	if _, ok := users[user_id]; !ok {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Error{Msg: "Нет пользователей"}
		message, _ := json.Marshal(msg)
		w.Write(message)
	}
	users[user_id].Email = tmp.Email
	users[user_id].First_name = tmp.First_name
	users[user_id].Last_name = tmp.Last_name
	users[user_id].Username = tmp.Username
}

//Logout ..
func Logout(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Cookie("sessionid")
	id2 := id.Value
	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   id2,
		Expires: time.Now(),
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	message, _ := json.Marshal(id2)
	w.Write(message)
}

//Leaders ..
func Leaders(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {
	numberOfPage := 0
	countOfString := 3
	canNext := true

	if r.URL.Query()["numPage"] != nil {
		numberOfPage, _ = strconv.Atoi(r.URL.Query()["numPage"][0])

	}

	//values2 := []interface{}{}
	values := make([]*models.User, 0, len(users))

	for _, value := range users {
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].Score > values[j].Score
	})

	// проверяем можно ли дальше листать
	if int(math.Ceil(float64(len(users))/float64(countOfString)))-1 < numberOfPage+1 {
		canNext = false
	}

	values = values[numberOfPage*countOfString : numberOfPage*countOfString+countOfString]
	//for _, value := range values {
	//	values2 = append(values2, value)
	//}
	//
	//values2 = append(values2, canNext)
	b := models.Leaders{}
	b.Users = values
	b.CanNext = canNext
	message, _ := json.Marshal(b)

	w.Header().Set("Content-Type", "application/json")
	w.Write(message)
}

//Upload ..
func Upload(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	file, handle, err := r.FormFile("image")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer file.Close()

	id, err := r.Cookie("sessionid")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	user_id := id.Value

	saveFile(users, w, file, user_id, handle)
}

func saveFile(users map[string]*models.User, w http.ResponseWriter, file multipart.File, user_id string, handle *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println(pwd)
	// src := pwd + "/uploads/" + users[user_id].Email + ".jpeg"
	src := pwd + "/uploads/" + users[user_id].Email + handle.Filename

	err = ioutil.WriteFile(src, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	users[user_id].Avatar = src

	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

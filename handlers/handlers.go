package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
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

	"2018_2_YetAnotherGame/models"
)

//SignUp ..
func SignUp(db *gorm.DB,ids map[string]string, users models.UsersMap, w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	tmp:=models.Session{}
	var user_id string
	db.Table("sessions").Select("id, email").Where("email = ?", user.Email).Scan(&tmp)
	user_id=tmp.ID

	if user_id!=""{
		fmt.Println("errrpr")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Error{Msg: "already exists"}
		message, _ := json.Marshal(msg)
		w.Write(message)
		return
	}

	id, _ := exec.Command("uuidgen").Output()

	stringID := string(id[:])
	stringID = strings.Trim(stringID, "\n")
	//users[stringID] = &user
	user.ID=stringID
	tmp.ID=stringID
	tmp.Email=user.Email
	db.Create(&tmp)
	//users.Store(stringID, &user)
	db.Create(&user)
	//ids[user.Email] = stringID

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
func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	cred := models.Auth{}
	tmp:=models.Session{}
	json.NewDecoder(r.Body).Decode(&cred)
	var user_id string
	db.Table("sessions").Select("id, email").Where("email = ?", cred.Email).Scan(&tmp)
user_id=tmp.ID
	if cred.Password == "" || cred.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Error{Msg: "Не указан E-Mail или пароль"}
		message, _ := json.Marshal(msg)
		w.Write(message)
		return
	}

	if user_id=="" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal("Неверный E-Mail")
		w.Write(message)
		return
	}
	var tmp2 models.User
	db.Where("id = ?", user_id).First(&tmp2)

	if tmp2.Password != cred.Password {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Error{Msg: "Неверный пароль"}
		message, _ := json.Marshal(msg)

		w.Write(message)
		return
	}


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
func Me(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := r.Cookie("sessionid")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id2 := id.Value
	tmp := models.Session{}
	db.Table("sessions").Select("id, email").Where("id = ?", id2).Scan(&tmp)
	//
	// db.Where("id = ?", id2).First(&tmp)
	if tmp.ID=="" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res:=models.User{}
	db.Where("id = ?", tmp.ID).First(&res)
//tmp,_:=users.Load(id2)
//	tmp.Score++

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	message, _ := json.Marshal(res)
	w.Write(message)
}

//Update ..
func Update(db *gorm.DB,  w http.ResponseWriter, r *http.Request) {
	tmp := models.User{}
	json.NewDecoder(r.Body).Decode(&tmp)
	id, _ := r.Cookie("sessionid")
	id2 := id.Value
	user_id := id2
	//fmt.Println(r.MatchString("peach"))
	ses:=models.Session{}


	db.Table("sessions").Select("id, email").Where("email = ?", user_id).Scan(&ses)

	if ses.ID=="" {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		msg := models.Error{Msg: "Нет пользователей"}
		message, _ := json.Marshal(msg)
		w.Write(message)
	}
	tmpuser:=models.User{}
	db.Table("users ").Select("id, email, first_name, last_name, username, password, score, avatar ").Where("id = ?", user_id).Scan(&tmpuser)

	db.Model(&tmpuser).Updates(models.User{Email:tmp.Email,First_name:tmp.First_name, Last_name:tmp.Last_name,Username:tmp.Username}).Where("id = ?", user_id)
	db.Save(&tmpuser)

}

//Logout ..
func Logout(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Cookie("sessionid")
	id2 := id.Value
	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   id2,
		Expires: time.Now(),
		Path: "/",

	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	message, _ := json.Marshal(id2)
	w.Write(message)
}

//Leaders ..
func Leaders(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	numberOfPage := 0
	countOfString := 3
	canNext := true

	if r.URL.Query()["numPage"] != nil {
		numberOfPage, _ = strconv.Atoi(r.URL.Query()["numPage"][0])

	}

	//values2 := []interface{}{}
	users := make([]*models.User, 0)
	query:="SELECT ID::text, email::text, first_name::text, last_name::text,username::text, score::integer FROM users;"

	rows,_ := db.Raw(query).Rows()




	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.ID,&user.Email,&user.First_name,&user.Last_name, &user.Username,&user.Score)
		if err != nil {		}

		users = append(users, user)
	}



	values := make([]*models.User, 0)

	for _, value := range users {
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].Score > values[j].Score
	})

	// проверяем можно ли дальше листать
	if int(math.Ceil(float64(len(values))/float64(countOfString)))-1 < numberOfPage+1 {
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
func Upload(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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

	saveFile(db, w, file, user_id, handle)
}

func saveFile(db *gorm.DB, w http.ResponseWriter, file multipart.File, user_id string, handle *multipart.FileHeader) {
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
	tmpuser:=models.User{}
//db.Table("users ").Select("id, email, first_name, last_name, username ").Where("id = ?", tmpuser).Scan(&tmpuser)
	src := pwd + "/uploads/" + tmpuser.Email + handle.Filename

	err = ioutil.WriteFile(src, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	//tmpuser.Avatar = src
	db.Table("users ").Select("id, email, first_name, last_name, username, password, score, avatar ").Where("id = ?", user_id).Scan(&tmpuser)
	//users.Store(user_id, tmpuser)
	db.Model(&tmpuser).Updates(models.User{Avatar:src}).Where("id = ?", user_id)
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

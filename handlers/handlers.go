package handlers

import (
	"encoding/json"
	"fmt"
	"goback/models"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func SignUp(ids map[string]string, users map[string]*models.User, w http.ResponseWriter, r *http.Request) {

	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	if _, ok := ids[user.Email]; ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal("already exists")
		w.Write(message)
	}

	id, _ := exec.Command("uuidgen").Output()

	stringId := string(id[:])
	stringId = strings.Trim(stringId, "\n")
	users[stringId] = &user
	ids[user.Email] = stringId

	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   stringId,
		Expires: time.Now().Add(60 * time.Minute),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)

	message, _ := json.Marshal(stringId)
	w.Write(message)

}

func Login(ids map[string]string, users map[string]*models.User, w http.ResponseWriter, r *http.Request) {
	cred := models.Auth{}
	json.NewDecoder(r.Body).Decode(&cred)
	user_id := ids[cred.Email]
	if cred.Password == "" || cred.Email == "" {

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal("Не указан E-Mail или пароль")
		w.Write(message)
	}

	if _, ok := users[user_id]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal("Неверный E-Mail")
		w.Write(message)
	}
	if users[user_id].Password != cred.Password {

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal("Неверный пароль")
		w.Write(message)
	}

	ids[cred.Email] = user_id
	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   user_id,
		Expires: time.Now().Add(60 * time.Minute),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)

	message, _ := json.Marshal(user_id)
	w.Write(message)

}

func Me(users map[string]*models.User, avatars map[string]string, w http.ResponseWriter, r *http.Request) {
	body := []interface{}{}

	id, _ := r.Cookie("sessionid")
	id2 := id.Value

	if _, ok := users[id2]; !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	body = append(body, users[id2])
	users[id2].Score += 1

	if src, ok := avatars[users[id2].email]; ok {
		body = append(body, src)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	message, _ := json.Marshal(body)
	w.Write(message)
}

func Update(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {

	tmp := models.User{}
	json.NewDecoder(r.Body).Decode(&tmp)
	id, _ := r.Cookie("sessionid")
	id2 := id.Value
	user_id := id2
	//fmt.Println(r.MatchString("peach"))

	if _, ok := users[user_id]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		message, _ := json.Marshal("no users")
		w.Write(message)
	}
	users[user_id].Email = tmp.Email
	users[user_id].First_name = tmp.First_name
	users[user_id].Last_name = tmp.Last_name
	users[user_id].Username = tmp.Username

}

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

func Leaders(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {

	limit := 6
	offset := 0

	if r.URL.Query()["limit"] != nil {
		limit, _ = strconv.Atoi(r.URL.Query()["limit"][0])

	}

	if r.URL.Query()["offset"] != nil {
		offset, _ = strconv.Atoi(r.URL.Query()["offset"][0])

	}

	values2 := []interface{}{}
	//fmt.Println(limit,offset)
	values := make([]*models.User, 0, len(users))

	for _, value := range users {
		values = append(values, value)
	}
	l := len(values)
	sort.Slice(values, func(i, j int) bool {
		return values[i].Score > values[j].Score
	})
	values = values[offset : limit+offset]
	for _, value := range values {
		values2 = append(values2, value)
	}
	values2 = append(values2, l)
	message, _ := json.Marshal(values2)

	w.Header().Set("Content-Type", "application/json")
	w.Write(message)

}

func Upload(avatars map[string]string, w http.ResponseWriter, r *http.Request) {
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

	r.ParseMultipartForm(0)
	email := r.FormValue("email")
	fmt.Println(email)

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg", "image/png":
		// saveFile(w, file, handle)
		saveFile(w, file, email, handle, avatars)
	default:
		jsonResponse(w, http.StatusBadRequest, "The format file is not valid.")
	}
}

func saveFile(w http.ResponseWriter, file multipart.File, email string, handle *multipart.FileHeader, avatars map[string]string) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	src := "../uploads/" + handle.Filename
	err = ioutil.WriteFile(src, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	avatars[email] = src

	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

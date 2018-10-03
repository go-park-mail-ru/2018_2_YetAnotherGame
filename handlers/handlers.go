package handlers

import (
	"encoding/json"
	"goback/models"
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
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
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	message, _ := json.Marshal(stringId)
	w.Write(message)

}

func Login(ids map[string]string, users map[string]*models.User, w http.ResponseWriter, r *http.Request) {
	cred := models.Auth{}
	json.NewDecoder(r.Body).Decode(&cred)
	user_id := ids[cred.Email]
	if cred.Password == "" || cred.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		message, _ := json.Marshal("Неверный пароль")
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

func Me(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {

	id, _ := r.Cookie("sessionid")
	id2 := id.Value
	if _, ok := users[id2]; !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	users[id2].Score += 1
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	message, _ := json.Marshal(users[id2])
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
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

	if r.URL.Query()["limit"] != nil {
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

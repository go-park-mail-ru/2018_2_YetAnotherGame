package handlers

import (
	"encoding/json"
	"goback/models"
	"net/http"
	"os/exec"
	"sort"
	"strconv"
	"time"
)



func SignUp(ids map[string]string, users map[string]*models.User, w http.ResponseWriter, r *http.Request) {

	//password := r.FormValue("password")
	//email := r.FormValue("email")
	//username := r.FormValue("username")
	//first_name := r.FormValue("first_name")
	//last_name := r.FormValue("last_name")
	user:=models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	if _, ok := ids[user.Email]; ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal("already exists")
		w.Write(b)
	}

	id, _ := exec.Command("uuidgen").Output()
	//n := bytes.Index(id, []byte{0})
	stringId := string(id[:])
	users[stringId]=&user
	ids[user.Email]=stringId

	cookie:= &http.Cookie{
		Name:"sessionid",
		Value:stringId,
		Expires:time.Now().Add(60*time.Minute),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)

	b, _ := json.Marshal(stringId)
	w.Write(b)
}




func Login(ids map[string]string, users map[string]*models.User, w http.ResponseWriter, r *http.Request) {

	password := r.FormValue("password")
	email := r.FormValue("email")
	user_id := ids[email]
	if password == "" || email == "" {

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal("Не указан E-Mail или пароль")
		w.Write(b)
	}

	if _, ok := users[user_id]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal("Неверный E-Mail")
		w.Write(b)
	}
	if users[user_id].Password != password {

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal("Неверный пароль")
		w.Write(b)
	}

	ids[email] = user_id
	cookie := &http.Cookie{
		Name:    "sessionid",
		Value:   user_id,
		Expires: time.Now().Add(60 * time.Minute),
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)

	b, _ := json.Marshal(user_id)
	w.Write(b)

}


func Me(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {

	id,_:= r.Cookie("sessionid")
	id2:=id.Value
	if _, ok := users[id2]; !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	users[id2].Score+=1
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(users[id2])
	w.Write(b)
}



func Leaders(users map[string]*models.User, w http.ResponseWriter, r *http.Request) {


	//fmt.Println(c)
limit:=6
offset:=0

if r.URL.Query()["limit"]!=nil{
	limit, _=strconv.Atoi(r.URL.Query()["limit"][0])

}

	if r.URL.Query()["offset"]!=nil{
		offset, _=strconv.Atoi(r.URL.Query()["offset"][0])

	}

	values2 := []interface{}{}
	//fmt.Println(limit,offset)
	values := make([]*models.User, 0, len(users))

	for  _, value := range users {
		values = append(values, value)
	}
	l:=len(values)
	//fmt.Println(eatureVector)
	sort.Slice(values, func(i, j int) bool {
		return values[i].Score > values[j].Score
	})
	values=values[offset:limit+offset]
	for  _, value := range values {
		values2 = append(values2, value)
	}
	values2 = append(values2, l)
	//fmt.Println(values)
	b, _:= json.Marshal(values2)
	//fmt.Println(b)

	//fmt.Println(b)
	//t, _:= json.Marshal(len(values))

	//fmt.Println(len(values))
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	//w.Write(t)

}
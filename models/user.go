package models

import "sync"

type User struct {
	ID        string `gorm:"primary_key"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Score      int    `json:"score"`
	Avatar     string `json:"avatar"`
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Error struct {
	Msg string `json:"Msg"`
}

type Leaders struct {
	Users   []*User
	CanNext bool `json:"CanNext"`
}


type UsersMap struct {
	mx sync.Mutex
	M map[string]*User
	Size int
}

func (c *UsersMap) Const()  {
	c.Size=0
	c.M=make(map[string]*User, 0)
}


func (c *UsersMap) Load(key string) (*User, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.M[key]
	return val, ok
}

func (c *UsersMap) Store(key string, value *User) {
	c.Size++
	c.mx.Lock()
	defer c.mx.Unlock()
	c.M[key] = value
}
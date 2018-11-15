package middlewares

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

/*	Берем данные которые приходят JSON'ом в теле ответа и парсим их в
 *	VKRosonseData добавляем к этому email и все эти данные пушим в проверяя нет ли такого юзера
 *	!!!Проблема: в вк может и не быть email'a, и тогда поле email в базе будет пустым, а оно ключевое
 *	Вариант решения: просить пользователя самому его вводить
 *
 */

const (
	ClientID     = "6744106"
	ClientKey    = "PLPNUPHSBYeveTNLb87w"
	ClientSecret = "5d968a9e5d968a9e5d968a9e3f5df062b455d965d968a9e067401d51f1b0c1062de795b"
	RedirectURL  = "http://127.0.0.1:8000"
	API_URL      = "https://api.vk.com/method/users.get?fields=photo_100&access_token=%s&v=5.52"
)

type VKResponseData struct {
	Response []struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Photo     string `json:"photo_100"`
	}
	Email string
}

func OauthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		code := r.FormValue("code")
		conf := oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  RedirectURL,
			Endpoint:     vk.Endpoint,
		}

		token, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Printf("Bad Exchange: %v", err)
			return
		}
		email := token.Extra("code").(string)

		client := conf.Client(ctx, token)

		resp, err := client.Get(fmt.Sprintf(API_URL, token.AccessToken))
		if err != nil {
			log.Printf("Bad resp: %v", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Bad read body: %v", err)
			return
		}

		data := &VKResponseData{}
		json.Unmarshal(body, data)
		data.Email = email //Теперь из Имени и Фамилии можно сделать Юзернейм и не будет только пароля, но можно добавить токен для валидации

		handler.ServeHTTP(w, r)
	})
}

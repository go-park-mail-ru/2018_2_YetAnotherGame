package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"2018_2_YetAnotherGame/ApiMS/resources/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

func (env *Environment) AvatarHandle(w http.ResponseWriter, r *http.Request) {
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

	aws_access_key_id := "d2ThzGBFkuZpxGezLQCFv2"
	aws_secret_access_key := "3PabWzRVmsAMog135pA85wmRdxPNMAMaRrMg2mQH7Wvs"
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	_, err = creds.Get()
	if err != nil {
		log.Println("ERROR GET CREDS", err)
	}

	cfg := aws.NewConfig().WithRegion("ru-msk").WithCredentials(creds).WithEndpoint("http://hb.bizmrg.com")
	svc := s3.New(session.New(), cfg)

	UploadFile(svc, env.DB, file, handle, id.Value)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	msg := models.Message{"File uploaded successfully!"}
	message, err := easyjson.Marshal(msg)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(message)
}

func UploadFile(svc *s3.S3, db *gorm.DB, file multipart.File, handle *multipart.FileHeader, user_id string) {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Error(err)
		return
	}

	fileSend := bytes.NewReader(fileBytes)

	tmpUser := models.User{}
	db.Table("users").
		Select("id, email, first_name, last_name, username, password, score, avatar").
		Where("id = ?", user_id).
		Scan(&tmpUser)

	path := "/uploads/" + tmpUser.Email + "-" + handle.Filename

	params := &s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String("yetanothergamebucket"),
		Key:    aws.String(path),
		Body:   fileSend,
	}
	resp, err := svc.PutObject(params)

	if err != nil {
		log.Println("ERROR PUT OBJECT", err)
	}
	fmt.Println("DONE", awsutil.StringValue(resp))

	src := "https://hb.bizmrg.com/yetanothergamebucket" + path
	db.Model(&tmpUser).
		Updates(models.User{Avatar: src}).
		Where("id = ?", user_id)
}

// func (env *Environment) AvatarHandle(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}
//
// 	file, handle, err := r.FormFile("image")
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// 	defer file.Close()
//
// 	id, err := r.Cookie("sessionid")
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// 	fmt.Println(id.Value)
//
// 	saveFile(env.DB, w, file, id.Value, handle)
// }
//
// func saveFile(db *gorm.DB, w http.ResponseWriter, file multipart.File, user_id string, handle *multipart.FileHeader) {
// 	data, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
//
// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		logrus.Error(err)
// 		os.Exit(1)
// 	}
//
// 	tmpuser := models.User{}
// 	db.Table("users ").Select("id, email, first_name, last_name, username, password, score, avatar ").Where("id = ?", user_id).Scan(&tmpuser)
// 	src := pwd + "/uploads/" + tmpuser.Email + handle.Filename
//
// 	err = ioutil.WriteFile(src, data, 0666)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// 	db.Model(&tmpuser).Updates(models.User{Avatar: src}).Where("id = ?", user_id)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	msg := models.Message{"File uploaded successfully!"}
// 	message, err := easyjson.Marshal(msg)
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	w.Write(message)
// }

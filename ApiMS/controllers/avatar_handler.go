package controllers

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"2018_2_YetAnotherGame/ApiMS/infostructures/functions"
	"2018_2_YetAnotherGame/ApiMS/resources/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (env *Environment) AvatarHandle(w http.ResponseWriter, r *http.Request) {
	id, err := r.Cookie("sessionid")
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(id.Value)

	file, handle, err := r.FormFile("image")
	if err != nil {
		logrus.Error(err)
		return
	}
	defer file.Close()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region:   aws.String("ru-msk"),
			Endpoint: aws.String("http://hb.bizmrg.com"),
		},
	}))
	svc := s3.New(sess)

	err = UploadFile(svc, env.DB, file, handle, id.Value)
	if err != nil {
		logrus.Error(err)
		functions.SendStatus("Cannot upload file", w, 500)
		return
	}

	functions.SendStatus("File uploaded successfully", w, 201)
}

func UploadFile(svc *s3.S3, db *gorm.DB, file multipart.File, handle *multipart.FileHeader, user_id string) error {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
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
		return err
	}
	logrus.Info(awsutil.StringValue(resp))

	src := "https://hb.bizmrg.com/yetanothergamebucket" + path
	db.Model(&tmpUser).
		Updates(models.User{Avatar: src}).
		Where("id = ?", user_id)

	return nil
}

package functions

import (
	"2018_2_YetAnotherGame/resources/models"
	"github.com/mailru/easyjson"
	"net/http"
)

func SendStatus(mes string, w http.ResponseWriter, header int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
	msg := models.Error{Msg: mes}
	message, err := easyjson.Marshal(msg)
	if err != nil {
		return err
	}
	w.Write(message)
	return nil
}


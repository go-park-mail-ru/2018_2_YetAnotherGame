package functions

import (
	"github.com/go-park-mail-ru/2018_2_YetAnotherGame/ApiMS/resources/models"
	"net/http"

	"github.com/mailru/easyjson"
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

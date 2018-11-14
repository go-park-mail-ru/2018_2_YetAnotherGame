package functions

import (
	"2018_2_YetAnotherGame/domain/models"
	"encoding/json"
	"net/http"
)

func BadRequest(mes string, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	msg := models.Error{Msg: mes}
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	w.Write(message)
	return nil
}

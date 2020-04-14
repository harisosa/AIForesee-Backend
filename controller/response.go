package controller

import (
	"encoding/json"
	"net/http"
	"sms-aiforesee-be/models"
)

func ResponseHandler(w http.ResponseWriter, data interface{}, message string, httpstatus int) {

	response := make(map[string]interface{})
	sucess := false
	if httpstatus == 200 {
		sucess = true
	}
	response = models.Message(sucess, message)
	response["data"] = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatus)

	json.NewEncoder(w).Encode(response)
}

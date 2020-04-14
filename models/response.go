package models

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Response(w http.ResponseWriter, data map[string]interface{}) {
	json.NewEncoder(w).Encode(data)
}

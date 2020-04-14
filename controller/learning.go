package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sms-aiforesee-be/database"
	"sms-aiforesee-be/models"
	"strings"
)

func LearningAScore(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var req map[string]interface{}
	err := decoder.Decode(&req)

	if err != nil {
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	if req["api_key"] == nil {
		ResponseHandler(w, nil, "api_key Requred", http.StatusInternalServerError)
		return
	}
	var apiUrl string
	if req["model_name"] == nil {
		ResponseHandler(w, nil, "model_name Requred", http.StatusInternalServerError)
		return
	}
	modelName := strings.ToLower(strings.TrimSpace(req["model_name"].(string)))

	var modelID models.APIType
	switch modelName {
	case "ascore":
		apiUrl += "http://localhost:2011/ascore_aiforesee"
		modelID = models.AScore
		break
	case "bscore":
		apiUrl += "http://localhost:2020/bscore_aiforesee"
		modelID = models.BScore
		break
	default:
		ResponseHandler(w, nil, "model name not found", http.StatusInternalServerError)
		return

	}

	db, err := database.ConnectDB()
	if err != nil {
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the existing entry present in the database for the given username
	result := tx.QueryRow("select * from users where api_key=$1", req["api_key"])

	var us models.User
	err = result.Scan(&us.ID,
		&us.Username,
		&us.Password,
		&us.Email,
		&us.ApiKey,
		&us.Created,
	)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			tx.Rollback()
			ResponseHandler(w, nil, "Api Key is not valid", http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		tx.Rollback()
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		tx.Rollback()
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	request, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		tx.Rollback()
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tx.Rollback()
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	var res map[string]interface{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		tx.Rollback()
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		tx.Rollback()
		ResponseHandler(w, res, "Error", resp.StatusCode)
		return
	}
	err = SaveEvent(tx, us.ID, req, res, modelID)
	if err != nil {
		tx.Rollback()
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()
	ResponseHandler(w, res, "", http.StatusOK)
	return
}

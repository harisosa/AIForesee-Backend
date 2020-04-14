package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sms-aiforesee-be/database"
	"sms-aiforesee-be/models"
	"time"

	"github.com/google/uuid"
)

func SaveEvent(tx *sql.Tx, usID string, req map[string]interface{}, res map[string]interface{}, types models.APIType) error {
	var evt models.Event
	evt.ID = uuid.New().String()
	evt.UserID = usID
	evt.Date = time.Now()
	request, _ := json.Marshal(req)
	response, _ := json.Marshal(res)
	evt.ApiType = types

	stmt, err := tx.Prepare(`INSERT INTO 
	events(id,user_id,date,request,response,api_type) 
	VALUES ($1, $2, $3, $4,$5,$6)`)
	if err == nil {
		_, err := stmt.Exec(evt.ID, evt.UserID, evt.Date, request, response, evt.ApiType)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func GetEventListByUserID(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]

	if !ok || len(ids[0]) < 1 {
		ResponseHandler(w, nil, "ID is Missing", http.StatusBadRequest)
		return
	}
	id := ids[0]

	db, err := database.ConnectDB()
	if err != nil {
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var evt models.Event
	var arr_evt []models.Event
	rows, err := db.Query("SELECT * FROM events WHERE user_id=$1", id)
	if err != nil {
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	var strReq []byte
	var strRes []byte
	for rows.Next() {
		if err := rows.Scan(&evt.ID, &evt.UserID, &evt.Date, &strReq, &strRes, &evt.ApiType); err != nil {
			ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
			return
		}
		json.Unmarshal([]byte(strReq), &evt.Request)
		json.Unmarshal([]byte(strRes), &evt.Response)
		arr_evt = append(arr_evt, evt)

	}
	ResponseHandler(w, arr_evt, "Success", http.StatusOK)
	return
}

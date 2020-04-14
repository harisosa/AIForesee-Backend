package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sms-aiforesee-be/database"
	"sms-aiforesee-be/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var us models.User
	err := decoder.Decode(&us)
	if err != nil {
		ResponseHandler(w, nil, err.Error(), http.StatusBadRequest)
		return
	}

	isExist, err := QueryUser(us.Username)
	if err != nil {
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
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
	if !isExist {
		us.ID = uuid.New().String()
		us.Created = time.Now()
		token, err := GenerateKey()
		if err != nil {
			tx.Rollback()
			ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
			return
		}
		us.ApiKey = token
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
		us.Password = string(hashedPassword)
		stmt, err := tx.Prepare(`INSERT INTO 
								users(id,username,password,email,api_key,created) 
								VALUES ($1, $2, $3, $4,$5,$6)`)
		if err == nil {
			_, err := stmt.Exec(us.ID, us.Username, us.Password, us.Email, us.ApiKey, us.Created)
			if err != nil {
				tx.Rollback()
				ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			tx.Rollback()
			ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
			return
		}
		tx.Commit()
		ResponseHandler(w, us, "Success", http.StatusOK)
	} else {
		ResponseHandler(w, nil, "Username Already Exist", http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		ResponseHandler(w, nil, err.Error(), http.StatusBadRequest)
		return
	}
	db, err := database.ConnectDB()
	if err != nil {
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the existing entry present in the database for the given username
	result := db.QueryRow("select * from users where username=$1  OR email=$1", creds.Username)
	if err != nil {
		// If there is an issue with the database, return a 500 error
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}
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
			ResponseHandler(w, nil, "Username or email not found", http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		ResponseHandler(w, nil, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(us.Password), []byte(creds.Password)); err != nil {
		// If there is an issue with the database, return a 500 error
		ResponseHandler(w, nil, "Wrong Password", http.StatusInternalServerError)
		return
	}
	ResponseHandler(w, us, "Login Sucessful ", http.StatusOK)
	return

}

func QueryUser(username string) (bool, error) {
	var users models.User
	db, err := database.ConnectDB()
	defer db.Close()
	err = db.QueryRow("SELECT ID FROM users WHERE username = $1", username).Scan(&users.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			return false, err
		}

		return false, nil
	}

	return true, err
}

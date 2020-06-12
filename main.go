package main

import (
	"fmt"
	"net/http"
	controller "sms-aiforesee-be/controller"
	"sms-aiforesee-be/migrations"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	err := migrations.Migrate()
	if err != nil {
		panic(err)
	}
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Origin", "Authorization", "X-Requested-With", "Content-Type", "Accept"}),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST", "PUT"}),
		//handlers.AllowCredentials(),
	)
	router := mux.NewRouter()
	router.Host("http://localhost/")

	corsMw := mux.CORSMethodMiddleware(router)

	//auth API
	authrouter := router.PathPrefix("/api/auth").Subrouter()
	authrouter.HandleFunc("/register", controller.Register).Methods("POST")
	authrouter.HandleFunc("/login", controller.Login).Methods("POST")

	learning := router.PathPrefix("/api/sms").Subrouter()
	learning.HandleFunc("/score", controller.LearningAScore).Methods("POST")

	event := router.PathPrefix("/api/event").Subrouter()
	event.HandleFunc("/getall", controller.GetEventListByUserID).Methods("GET")

	router.Use(corsMw)

	fmt.Println("Connected to SMS-AIforesee API")
	fmt.Println(http.ListenAndServe(":80", cors(router)))
}

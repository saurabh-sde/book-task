package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/saurabh-sde/library-task-go/controller"
	"github.com/saurabh-sde/library-task-go/middleware"
)

func init() {
	godotenv.Load()
}

func main() {
	r := mux.NewRouter()

	// login route
	r.HandleFunc("/login", controller.HandleLogin)

	// set auth middleware sub-router
	privateRouter := r.PathPrefix("/").Subrouter()
	privateRouter.Use(middleware.AuthMiddleware)

	// home route
	privateRouter.HandleFunc("/home", controller.HandleHome).Methods("GET")

	// book route
	privateRouter.HandleFunc("/addBook", controller.HandleAddBook).Methods("POST")
	privateRouter.HandleFunc("/deleteBook", controller.HandleDeleteBook).Methods("DELETE")

	// add loggin middleware
	r.Use(middleware.Logging)

	fmt.Println("Starting Local Server: http://localhost:8080")
	// start local server
	http.ListenAndServe(":8080", r)
}

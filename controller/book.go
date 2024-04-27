package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/saurabh-sde/library-task-go/service"
	"github.com/saurabh-sde/library-task-go/utility"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	userType := utility.UserTypeFromContext(r.Context())

	// default user type is not stated so set it to regular
	if userType == "" {
		userType = "regular"
	}

	// fetch books
	books, err := service.GetBooksByUserType(userType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func HandleAddBook(w http.ResponseWriter, r *http.Request) {
	userType := utility.UserTypeFromContext(r.Context())
	// return err for non-admin
	if userType != "admin" {
		http.Error(w, errors.New("permission denied").Error(), http.StatusUnauthorized)
		return
	}

	// req can be json which can be validated while unmarshel to model.Book{} struct
	// but considering query params to perform validation on input values

	// get data from URL query params
	params := r.URL.Query()
	utility.Print("HandleAddBook: ", params)
	bookName := params.Get("bookName")
	author := params.Get("author")
	pubYear := params.Get("publicationYear")
	// convert to int
	pubYearInt, err := strconv.Atoi(pubYear)
	if len(bookName) == 0 || len(author) == 0 || len(pubYear) == 0 || err != nil || pubYearInt <= 0 {
		http.Error(w, errors.New("invalid data").Error(), http.StatusBadRequest)
		return
	}

	// add book
	err = service.AddBook(bookName, author, pubYearInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	userType := utility.UserTypeFromContext(r.Context())
	if userType != "admin" {
		http.Error(w, errors.New("permission denied").Error(), http.StatusUnauthorized)
		return
	}
	params := r.URL.Query()
	bookName := params.Get("bookName")

	// delete book
	err := service.DeleteBook(bookName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/saurabh-sde/library-task-go/model"
	"github.com/saurabh-sde/library-task-go/utility"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	funcName := "HandleLogin"
	utility.Print(funcName)

	// read req body
	var userReq model.User
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		utility.Error("Err in reading: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(reqBody, &userReq)
	if err != nil {
		utility.Error("Err in Unmarshal: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utility.Print(userReq)

	// check email and password authentication with mock user table data
	// from there user type can also be known
	sampleUsersMap := utility.GetSampleUsersData()
	currUser, ok := sampleUsersMap[userReq.Email]
	if !ok {
		err = errors.New("user not found")
		utility.Error("No user found: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if currUser.Password != userReq.Password {
		err = errors.New("wrong password")
		utility.Error("Wrong Password: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// create jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    currUser.Email,
		"userType": currUser.UserType, // take from mocked user data
	})

	// get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(utility.GetSecret()))
	if err != nil {
		utility.Error("Err in SignedString: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"id":       currUser.Id,
		"email":    currUser.Email,
		"userType": currUser.UserType,
		"token":    tokenString,
	}
	jsonResp, err := json.Marshal(data)
	if err != nil {
		utility.Error("Err in Marshel: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Not needed as requirements to add into Cookie not given but it need to be sent into response
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// })

	// return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

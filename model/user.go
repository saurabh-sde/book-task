package model

type User struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserType string `json:"userType"`
	Password string `json:"password"`
}

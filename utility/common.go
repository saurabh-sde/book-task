package utility

import (
	"context"
	"fmt"
	"os"

	"github.com/saurabh-sde/library-task-go/model"
)

func Print(v ...interface{}) {
	logLvl := "INFO"
	fmt.Printf("%s: %+v\n", logLvl, v)
}

func Error(v ...interface{}) {
	logLvl := "ERROR"
	fmt.Printf("%s: %+v\n", logLvl, v)
}

func GetSecret() string { return os.Getenv("SECRET") }

func GetSampleUsersData() map[string]model.User {
	m := make(map[string]model.User)
	m["admin@test.com"] = model.User{Id: "adminId", Email: "admin@test.com", Password: "admin", UserType: "admin"}
	m["regular@test.com"] = model.User{Id: "regularId", Email: "regular@test.com", Password: "regular", UserType: "regular"}

	//TODO: add more user for testing login

	return m
}

type ContextKey string

const ContextUserType ContextKey = "userType"

func UserTypeFromContext(ctx context.Context) string {
	return ctx.Value(ContextUserType).(string)
}

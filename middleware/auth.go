package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/saurabh-sde/library-task-go/utility"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from request
		// assuming 'Authorization': 'Bearer <YOUR_TOKEN_HERE>' as standard
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer")
		tokenString := strings.TrimSpace(splitToken[1])

		// Check if user is authenticated by parsing token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// validate the alg for SigningMethodHS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(utility.GetSecret()), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		utility.Print("Successfully Validated: ", claims)

		// set user type into context
		// user id is not needed to fetch user else it can be added
		ctx := context.WithValue(r.Context(), utility.ContextUserType, claims["userType"])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

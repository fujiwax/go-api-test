package app

import (
	"context"
	"fmt"
	"go-contacts/models"
	u "lens/utils"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuth := []string{"/api/user/new", "/api/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                              //current request path

		// check if request does not need authenticaiton, serve the request if it doesn't need it
		for _, value := range noAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") // grab the token from the Authorization

		if tokenHeader == "" {
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") // `Bearer {token-body}` format
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //token body
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid {
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-type", "application/json")
			u.Respond(w, response)
			return
		}

		// everything went well, proceed with the request and set the caller to the user retrieved from parsed token
		fmt.Strintf("User %", tk.Username) // get user name
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

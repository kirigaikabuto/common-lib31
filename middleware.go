package common

import (
	"context"
	"encoding/json"
	"net/http"
)


type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}


type Middleware interface {
	LoginMiddleware(fn http.Handler) http.Handler
}

type middleware struct {
	redisStore *RedisStore
}

func NewMiddleware(rS *RedisStore) Middleware {
	return &middleware{redisStore: rS}
}

func (m *middleware) LoginMiddleware(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if key == "" {
			respondJSON(w, http.StatusInternalServerError, &HttpError{"For access needed authorization header", http.StatusInternalServerError})
			return
		}
		if key != "" {
			user := &User{}
			err := m.redisStore.GetValue(key, &user)
			if err != nil {
				errorMessage := err.Error()
				if err.Error() == "redis: nil" {
					errorMessage = "Your access key is expired"
				}
				respondJSON(w, http.StatusInternalServerError, HttpError{
					Message:    errorMessage,
					StatusCode: http.StatusInternalServerError,
				})
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", user.Id)
			r = r.WithContext(ctx)
		}
		fn.ServeHTTP(w, r)
	})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

type HttpError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

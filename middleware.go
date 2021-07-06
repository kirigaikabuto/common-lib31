package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Middleware interface {
	LoginMiddleware(fn http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	redisStore *RedisStore
}

func NewMiddleware(rS *RedisStore) Middleware {
	return &middleware{redisStore: rS}
}

type ErrorMessage struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (m *middleware) LoginMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if key == "" {
			respondJSON(w, http.StatusInternalServerError, &HttpError{"For access needed authorization header", http.StatusInternalServerError})
			return
		}
		if key != "" {
			userId, err := m.redisStore.GetValue(key)
			if err != nil && strings.Contains(err.Error(), "redis: nil") {
				respondJSON(w, http.StatusBadRequest, &ErrorMessage{
					Message: "Your token is expired",
					Status:  http.StatusBadRequest,
				})
				return
			} else if err != nil {
				respondJSON(w, http.StatusBadRequest, &ErrorMessage{
					Message: err.Error(),
					Status:  http.StatusBadRequest,
				})
				return
			}
			fmt.Println(userId)
			ctx := context.WithValue(r.Context(), "user_id", userId)
			r = r.WithContext(ctx)
		}
		fn.ServeHTTP(w, r)
	}
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

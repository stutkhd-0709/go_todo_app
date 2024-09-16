package handler

import (
	"fmt"
	"github.com/stutkhd-0709/go_todo_app/auth"
	"log"
	"net/http"
)

// AdminMiddleware はadminロールか判定する
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("AdminMiddleware")
		if !auth.IsAdmin(r.Context()) {
			RespondJSON(r.Context(), w, ErrorResponse{
				Message: "not admin",
			}, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware はミドルウェアパターンと合わせたインタフェースになっている
func AuthMiddleware(j *auth.JWTer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, err := j.FillContext(r)
			if err != nil {
				RespondJSON(r.Context(), w, ErrorResponse{
					Message: "not find auth info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)
				return
			}
			x, ok := auth.GetUserID(req.Context())
			if !ok {
				fmt.Println("not ok")
			}
			fmt.Println(x)
			next.ServeHTTP(w, req)
		})
	}
}

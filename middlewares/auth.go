package middlewares

import (
	"fmt"
	"net/http"
)

// jwt 权限验证中间件
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth before")
		//w.WriteHeader(403)
		//w.Write([]byte("权限不足"))
		//return
		next.ServeHTTP(w, r)
		fmt.Println("auth after")
	})
}

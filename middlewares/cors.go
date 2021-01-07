package middlewares

import (
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Origin,Content-Type,Content-Length,Authorization,Token,Cookie")
		w.Header().Add("Access-Control-Allow-Credentials", "true") //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,HEAD") //允许请求方法
		next.ServeHTTP(w, r)
	})
}
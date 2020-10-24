package middlewares

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)
// 请求日志中间件
func LoggerMiddleware(next http.Handler) http.Handler  {
	loggerPath := "request.log"
	logger := logrus.New()
	//logger.SetFormatter(&logrus.JSONFormatter{})
	logFile,_ := os.OpenFile(loggerPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	logger.SetOutput(logFile)
	logger.SetLevel(logrus.InfoLevel)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		next.ServeHTTP(w, r)
		// 响应时间
		execTime := time.Now().Sub(timeStart)
		requestMethod := r.Method
		remoteAddr := r.RemoteAddr
		requestURI := r.RequestURI

		logger.WithField("ip", remoteAddr).
			WithField("uri", requestURI).
			//WithField("status", httpStatus).
			WithField("method", requestMethod).
			WithField("duration", execTime).Info("")
	})
}
package middlewares

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)
// 请求日志中间件
func LoggerMiddleware(next http.Handler) http.Handler  {
	basePath := "log/"
	logger := logrus.New()
	//logger.SetFormatter(&logrus.JSONFormatter{})
	//logFile,_ := os.OpenFile(loggerBasePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	writer, _ := rotatelogs.New(
		basePath + "request-%Y%m%d%H%M.log",
		rotatelogs.WithMaxAge(time.Duration(30 * 24)*time.Hour),
		// 请求日志保存180天
		rotatelogs.WithRotationTime(time.Duration(180 * 24)*time.Hour),
	)
	logger.SetOutput(writer)
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
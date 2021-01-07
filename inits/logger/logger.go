package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/showiot/camera/inits/config"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)
var logger *logrus.Logger

func Init() error  {
	 logCfg:= config.Conf.LoggerConfig
	logger = logrus.New()
	//logger.SetFormatter(&logrus.JSONFormatter{})
	//logFile,err := os.OpenFile(loggerPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	//if err !=nil {
	//	return err
	//}
	writer, _ := rotatelogs.New(
		"log/%Y%m%d%H%M.log",
		//rotatelogs.WithLinkName(loggerPath),
		rotatelogs.WithMaxAge(time.Duration(logCfg.MaxAge * 24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(logCfg.RotationTime * 24)*time.Hour),
	)
	logger.SetOutput(writer)
	logger.SetLevel(logrus.InfoLevel)

	return nil
}
func GetLogger() *logrus.Entry {
	fileName, funcName, lineNo :=Trace(2)
	entry := logger.WithField("file", fileName).WithField("func", funcName).WithField("lineNo", lineNo)
	return entry
}
func Trace(skip int) (fileName string, funcName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}
	return
}
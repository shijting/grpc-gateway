package logger

import (
	"github.com/showiot/camera/inits/config"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)
var logger *logrus.Logger

func Init() error  {
	loggerPath := config.Conf.LoggerConfig.LoggerFile
	log.Println("logger file path:", loggerPath)
	logger = logrus.New()
	//logger.SetFormatter(&logrus.JSONFormatter{})
	logFile,err := os.OpenFile(loggerPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err !=nil {
		return err
	}
	logger.SetOutput(logFile)

	logger.SetLevel(logrus.InfoLevel)
	return nil
}
func GetLogger() *logrus.Logger {
	return logger
}
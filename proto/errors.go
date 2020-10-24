package proto

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Errorg(code codes.Code, errCode Error_Code, errMsgs ...string) error {
	var errMsg string
	if len(errMsgs) > 0 {
		errMsg = errMsgs[0]
	}
	if errMsg == "" {
		errMsg = errCode.String()
	}
	return status.New(code, errMsg).Err()
}

func Errorf(errCode Error_Code, errMsgs ...string) *Error {
	var errMsg string
	if len(errMsgs) > 0 {
		errMsg = errMsgs[0]
	}
	if errMsg == "" {
		errMsg = viper.GetString("errors." + errCode.String())
	}
	return &Error{
		Code:    errCode,
		Message: errMsg,
	}
}

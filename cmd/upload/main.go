package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/showiot/camera/inits/config"
	"github.com/showiot/camera/pkg/v1/upload"
	"github.com/showiot/camera/proto"
	"google.golang.org/grpc/codes"
	"log"
	"mime/multipart"
	"net/http"
	"path"
)

const maxImageSize = 1 << 20

type IUpload interface {
	Save(file multipart.File, size int64, ext string)( string,  error)
}


type errorResp struct {
	Code codes.Code `json:"code"`
	Msg  string     `json:"msg"`
}

var configPath = ""

func init()  {
	flag.StringVar(&configPath, "config_path", "configs/config.yaml", "-config_path \"configs/config.yaml\"")
}

func setupRoutes() {
	flag.Parse()
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8003", nil)
}

func main() {
	if err := config.Init(configPath); err != nil {
		log.Fatal(err)
	}
	setupRoutes()

}
func uploadFile(w http.ResponseWriter, r *http.Request) {
	var respByte []byte
	// TODO: 权限校验
	_ = r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File ", err)
		respByte, _ = marshalResp(codes.Internal, proto.Error_ERR_INTERNAL_SERVER.String())
		w.Write(respByte)
		return
	}
	defer file.Close()
	size := handler.Size
	ext := path.Ext(handler.Filename)
	diskStore := upload.NewDiskFileStore("uploads")
	filePath, err := diskStore.Save(file, size, ext)
	if err !=nil {
		fmt.Println(err)
		if err == upload.ErrFileSizeLimit {
			respByte, _ = marshalResp(codes.OutOfRange, proto.Error_ERR_UPLOAD_SIZE_LIMIT.String())
		} else {
			respByte, _ = marshalResp(codes.InvalidArgument, proto.Error_ERR_OPERATION_FAILED.String())
		}
		w.Write(respByte)
	}
	ossPath, err := upload.NewOssFileStore().Save(filePath, size, ext)
	if err !=nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(ossPath))
}
func marshalResp(code codes.Code , msg string) ([]byte, error) {
	resp := &errorResp{
		Code: code,
		Msg: msg,
	}
	respByte, err :=json.Marshal(resp)
	return respByte, err
}

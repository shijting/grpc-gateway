package upload

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/showiot/camera/inits/config"
	"io"
	"mime/multipart"
	"os"
	"time"
)

var ErrFileSizeLimit = errors.New("out of the file limit")

type diskFileStore struct {
	FileFolder string
}
// 本地存储
func NewDiskFileStore(fileFolder string) *diskFileStore {
	return &diskFileStore{FileFolder: fileFolder}
}
func (this *diskFileStore)Save(file multipart.File, size int64, ext string)(savePath string, err error)  {
	limit := config.Conf.UploadConfig.Size
	if limit < int(size) {
		err = ErrFileSizeLimit
		return
	}
	_uuid, err := uuid.NewRandom()
	fileFolder := fmt.Sprintf("%s/%s",this.FileFolder, time.Now().Format("200601"))
	os.MkdirAll(fileFolder, 0666)
	savePath = fmt.Sprintf("%s/%s%s",fileFolder, _uuid, ext)
	_saveFile, err := os.OpenFile(savePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer _saveFile.Close()
	if _, err = io.Copy(_saveFile, file); err != nil {
		return
	}
	return
}

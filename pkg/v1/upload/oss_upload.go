package upload

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/showiot/camera/inits/config"
	"io"
	"math"
	"os"
)

type ossFileStore struct {}

func NewOssFileStore() *ossFileStore {
	return &ossFileStore{}
}

func (this *ossFileStore)Save(locaFilename string, size int64, ext string)(savePath string, err error)  {
	ossConfig := config.Conf.OssConfig
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	bucketName := ossConfig.BucketName
	_uuid, err := uuid.NewRandom()
	savePath = fmt.Sprintf("%s%s", _uuid, ext)
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	chunks, err := oss.SplitFileByPartNum(locaFilename, partNum(size))
	var fd *os.File
	fd, err = os.Open(locaFilename)
	defer fd.Close()

	// 指定存储类型为标准存储。
	storageType := oss.ObjectStorageClass(oss.StorageStandard)

	// 步骤1：初始化一个分片上传事件，并指定存储类型为标准存储。
	imur, err := bucket.InitiateMultipartUpload(savePath, storageType)
	// 步骤2：上传分片。
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		fd.Seek(chunk.Offset, io.SeekStart)
		// 调用UploadPart方法上传每个分片。
		part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		if err != nil {
			fmt.Println("Error:", err)
			return "", err
		}
		parts = append(parts, part)
	}
	// 指定Object的读写权限为公共读，默认为继承Bucket的读写权限。
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	// 步骤3：完成分片上传，指定文件读写权限为公共读。
	_, err = bucket.CompleteMultipartUpload(imur, parts, objectAcl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	//fmt.Println("cmur:", cmur)
	return
}
func partNum(size int64) int  {
	var partSize int64 = 1024 * 1024 * 3
	return int(math.Ceil(float64(size) / float64(partSize)))
}

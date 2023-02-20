package minio

import (
	"Tiktok/config"
	"Tiktok/pkg/log"
	"bytes"
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

// 初始化，与minio服务端建立连接
func init(){
	endpoint := config.MinioSetting.Host + ":" + config.MinioSetting.Port
	accessKeyID := config.MinioSetting.AccessKeyID
	secretAccessKey := config.MinioSetting.SecretAccessKey
	useSSL := false
	var err error
	client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Error(err.Error())
	}
	if err != nil {
		log.Error(err.Error())
	} else{
		log.Info("init minio success")
	}
}

// 创建桶，返回错误信息
func CreateBucket(ctx context.Context, bucketName string) error {
	location := "cn-north-1"
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		return err
	}
	s := "Create bucket " + bucketName + " successful"
	log.Info(s)
	return nil
}

// 检测桶是否存在
// 返回是否存在和错误信息
func DetectExist(ctx context.Context, bucketName string) (bool, error){
	found, err := client.BucketExists(ctx, bucketName)
	if err != nil{
		return false, err
	}
	return found, nil
}

// 上传文件
// 传入参数：桶名、文件路径、存储的名字、文件类型(video, cover)，返回错误信息
// 永久访问的url为127.0.0.1:9000/video/objectName , 封面把video换成cover
func UploadFile(bucketName string, data []byte, objectName string, contentType string) error {
	ctx := context.Background()
	reader := bytes.NewReader(data)
	found, err := client.BucketExists(ctx, bucketName)
	if err != nil{
		log.Error(err.Error())
		return err
	}
	if !found{
		err = CreateBucket(ctx, bucketName)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	_, err = client.PutObject(ctx, bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// 获取文件
// 传入桶名和存储的名字，返回url信息和错误信息
func GetFile(bucketName string, objectName string) (*url.URL, error){
	ctx := context.Background()
	reqParams := make(url.Values)
	presignedUrl, err := client.PresignedGetObject(ctx, bucketName, objectName, time.Hour, reqParams)
	if err !=nil {
		log.Error(err.Error())
		return nil, err
	}
	return presignedUrl, nil
}
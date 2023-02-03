package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

// 初始化，与minio服务端建立连接
func init(){
	endpoint := "127.0.0.1:9000"
	accessKeyID := "admin"
	secretAccessKey := "1766551219.qwe."
	useSSL := false
	var err error
	client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	found, err := DetectExist(context.Background(), "vedio")
	if err != nil {
		fmt.Println(err.Error())
	}
	if !found{
		err = CreateBucket("vedio")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// 创建桶，返回错误信息
func CreateBucket(bucketName string) error {
	ctx := context.Background()
	location := "cn-north-1"
	err := client.MakeBucket(ctx, "vedio", minio.MakeBucketOptions{Region: location})
	if err != nil {
		return err
	}
	fmt.Printf("Create bucket %s successful", bucketName)
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
// 传入参数：桶名、文件路径、存储的名字、文件类型，返回错误信息
func UploadFile(bucketName string, path string, objectName string, contentType string) error {
	ctx := context.Background()
	found, err := client.BucketExists(ctx, bucketName)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	if !found{
		err = CreateBucket(bucketName)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	info, err := client.FPutObject(ctx, bucketName, objectName, path, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(info.Size)
	return nil
}

// 获取文件
// 传入桶名和存储的名字，返回url信息和错误信息
func GetFile(bucketName string, objectName string) (*url.URL, error){
	ctx := context.Background()
	reqParams := make(url.Values)
	presignedUrl, err := client.PresignedGetObject(ctx, bucketName, objectName, time.Hour, reqParams)
	if err !=nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return presignedUrl, nil
}
package minio

import (
	"fmt"
	"testing"
)

func TestMinioUp(t *testing.T) {
	err := UploadFile("vedio", "../../img/Tiktok/test.png", "test", "application/png")
	if err != nil{
		t.Errorf(err.Error())
	}
}

func TestMinioGet(t *testing.T) {
	url, err := GetFile("video", "test")
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(url.String())
}

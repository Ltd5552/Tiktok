package minio

import (
	"fmt"
	"os"
	"testing"
)

func TestMinioUp(t *testing.T) {
	data, err := os.ReadFile("../../img/test.png")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = UploadFile("cover", data, "test", "cover")
	if err != nil{
		t.Errorf(err.Error())
	}
}

func TestMinioGet(t *testing.T) {
	url, err := GetFile("cover", "test")
	if err !=nil {
		t.Errorf(err.Error())
	}
	fmt.Println(url.String())
}

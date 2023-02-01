package jwt

import "testing"

func TestJwt(t *testing.T){
	token, err := CreateToken("123", "xiaoming")
	if err != nil{
		t.Errorf(err.Error())
	}
	id, err := ParseToken(token)
	if err !=nil{
		t.Errorf(err.Error())
	}else if id !=123 {
		t.Errorf("id is not true")
	}
}
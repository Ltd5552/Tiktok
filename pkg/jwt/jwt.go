package jwt

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	ID   string
	Name string
	jwt.RegisteredClaims
}

var MySecret = []byte("密令123")

func CreateToken(id string, name string) (string, error) {
	claim := Claims{
		ID:               id,
		Name:             name,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(MySecret)
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	}
}

func ParseToken(tokenstr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenstr, &Claims{}, Secret())
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return 0, errors.New("not a true token")
			} else{
				return 0, errors.New("unknow error")
			}
		}
	}
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		id, err := strconv.Atoi(claims.ID)
		if err!=nil {
			return 0, errors.New("ID is not int")
		}
		return id, nil
	}
	return 0, errors.New("couldn`t parse the token")
}

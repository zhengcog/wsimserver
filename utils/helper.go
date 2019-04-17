package utils

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"strconv"
	"strings"
	"time"
)

var (
	jwtsecret []byte = []byte("wd8qhn13JhZaklAd")
)

type MyCustomClaims struct {
	UserId int64
	jwt.StandardClaims
}

//生成json web token
func GenJWToken(userid int64) string {
	claims := MyCustomClaims{
		userid,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix()),
			ExpiresAt: int64(time.Now().Unix() + 3600*24*180), //有效期
			Issuer:    "zhengcog@gmail.com",                   //发行者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtsecret)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return ss
}

//从json web token中解析 得到UserId
func ParseJWTokenUserId(authString string) (userid int64, err error) {
	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		err = errors.New("AuthString invalid:" + authString)
		return 0, err
	}
	tokenString := kv[1]

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtsecret, err
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = errors.New("That‘s not even a token")
				return 0, err
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				err = errors.New("Token is either expired or not active yet")
				return 0, err
			} else {
				err = errors.New("Couldn‘t handle this token")
				return 0, err
			}
		} else {
			err = errors.New("Couldn‘t handle this token")
			return 0, err
		}
	}
	if !token.Valid {
		err = errors.New("Token invalid:" + tokenString)
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Couldn‘t Get UserId")
	}
	userid, _ = strconv.ParseInt(fmt.Sprintf("%.0f", claims["UserId"]), 10, 64)
	return userid, nil
}

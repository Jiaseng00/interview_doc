package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

func main() {
	// 用户登入生成JWT
	token, err := GenerateToken("example123")
	if err != nil {
		log.Fatal(err)
	}
	ParseToken(token)
}

func GenerateToken(username string) (string, error) {
	// 模拟用户请求登入Api，回车需要的资料
	// JWT Payload资料生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 放入需要的用户资料
		"foo":      "bar",
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	// 从.env文件，调用密钥来进行签名
	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tokenString)
	return tokenString, nil
}

func ParseToken(tokenString string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := []byte(os.Getenv("JWT_SECRET"))

		return secret, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["exp"])
	} else {
		fmt.Println(err)
	}
}

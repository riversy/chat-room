package users

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"os"
)

type User struct {
	Uuid string
	Name string
}

func NewUser(name string) *User {
	return &User{
		Uuid: uuid.NewString(),
		Name: name,
	}
}

func (u User) GetJwt() (string, error) {
	return GetJwtFromUuid(u.Uuid)
}

func GetJwtFromUuid(uuid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": uuid,
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func GetUserUuidFromJwt(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprint(claims["uuid"]), nil
	} else {
		return "", errors.New("provided JWT is not valid")
	}
}

func GetUserByUuid(uuid string) (*User, error) {
	// @todo add and implementation
	return nil, nil
}

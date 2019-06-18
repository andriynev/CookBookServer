package jwt_auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"

	"time"
)

var signingKey []byte

var keyFn = func(*jwt.Token) (interface{}, error) {
	return signingKey, nil
}

var parser = jwt.Parser{
	ValidMethods: []string{"HS256"},
}

func Setup(key string) {
	signingKey = []byte(key)
}

var ErrNotValidToken = errors.New("given access token is not valid")

func ParseToken(accessToken string) (claims *UserClaims, err error) {
	token, err := parser.ParseWithClaims(accessToken, &UserClaims{}, keyFn)
	if err != nil {
		return
	}

	if !token.Valid {
		err = ErrNotValidToken
		return
	}
	claims = token.Claims.(*UserClaims)
	return
}

type Privileges map[string]bool

func (p Privileges) Has(privilege string) bool {
	for privilegeName, privilegeValue := range p {
		if privilegeName == privilege && privilegeValue {
			return true
		}
	}

	return false
}

type UserClaims struct {
	jwt.StandardClaims
	Id         uint       `json:"id"`
	Privileges Privileges `json:"privileges"`
}

func (c UserClaims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}

	if c.Id == 0 {
		err = fmt.Errorf("user id cannot be empty")
		return err
	}
	return nil
}

func GetUserClaims(id uint, privileges map[string]bool) jwt.Claims {
	return &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(),
		},
		Id:         id,
		Privileges: privileges,
	}
}

func GetToken(claims jwt.Claims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(signingKey)
	if err != nil {
		return
	}

	return
}

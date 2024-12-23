package jwt

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	key []byte
}

func NewJWT(key string) *JWT {
	return &JWT{
		key: []byte(key),
	}
}

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrForbiden         = errors.New("forbiden")
	ErrTokenExpired     = errors.New("token is expired")
)

type UserClaims struct {
	ID       int64
	UserType core.UsersUserType
	jwt.StandardClaims
}

func (j JWT) GenUserToken(id int64, userType core.UsersUserType, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		ID:       id,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j JWT) VerifyToken(tokneString string, allowedUsers []core.UsersUserType) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokneString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, ErrInvalidSignature
	}

	claims := token.Claims.(*UserClaims)
	if !slices.Contains(allowedUsers, claims.UserType) {
		return nil, ErrForbiden
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrTokenExpired
	}

	return claims, nil
}

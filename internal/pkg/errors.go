package pkg

import (
	"errors"
	"fmt"
)

func PrintErr(wrapper, err error) {
	fmt.Printf("[ERROR] %v: [%v]\n", wrapper, err)
}

// transport errors
var ErrPayloadDecode = errors.New("can't decode payload")

// db erorrs
var (
	ErrDbInternal      = errors.New("db intenrnal error")
	ErrEmailDuplicate  = errors.New("email is already taken")
	ErrRetrievingUser  = errors.New("can't retrieve user")
	ErrNotFound        = errors.New("entity does not exist")
	ErrEntityDuplicate = errors.New("entity already exists")
)

// third party errors
var ErrTranslation = errors.New("couldn't translate text")

// service errors
var (
	ErrEncryptingPassword = errors.New("can't encrypt password")
	ErrWrongPassword      = errors.New("wrong password")
)

// jwt errors
var (
	ErrCreatingToken    = errors.New("error while creating token")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrForbiden         = errors.New("forbiden")
	ErrTokenExpired     = errors.New("token is expired")
)

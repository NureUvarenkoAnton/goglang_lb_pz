package service

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/jwt"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   iAuthRepo
	jwtHandler jwt.JWT
}

const tokenTimeToLive = time.Hour * 24 * 7

type iAuthRepo interface {
	CreateUser(ctx context.Context, arg core.CreateUserParams) error
	GetUserByEmail(ctx context.Context, email sql.NullString) (core.User, error)
	GetPetById(ctx context.Context, id int64) (core.Pet, error)
}

func NewAuthService(repo iAuthRepo, jwtHandler jwt.JWT) *AuthService {
	return &AuthService{
		userRepo:   repo,
		jwtHandler: jwtHandler,
	}
}

func (s AuthService) RegisterUser(ctx context.Context, user core.CreateUserParams, toHashPassword bool) (string, error) {
	u, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		pkg.PrintErr(pkg.ErrDbInternal, err)
		fmt.Printf("%T", err)
		return "", fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	if u.UserType == user.UserType {
		return "", pkg.ErrEmailDuplicate
	}

	if toHashPassword {
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password.String), bcrypt.DefaultCost)
		if err != nil {
			return "", fmt.Errorf("%w: [%w]", pkg.ErrEncryptingPassword, err)
		}

		user.Password.String = string(pass)
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			err := err.(*mysql.MySQLError)
			// duplicate entry
			if err.Number == 1062 {
				return "", pkg.ErrEmailDuplicate
			}

		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return "", pkg.ErrDbInternal
	}

	dbUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		pkg.PrintErr(pkg.ErrRetrievingUser, err)
		return "", pkg.ErrRetrievingUser
	}

	token, err := s.jwtHandler.GenUserToken(dbUser.ID, dbUser.UserType.UsersUserType, time.Now().Add(tokenTimeToLive))
	if err != nil {
		pkg.PrintErr(pkg.ErrCreatingToken, err)
		return "", pkg.ErrCreatingToken
	}

	return token, nil
}

func (s AuthService) Login(ctx context.Context, payload core.CreateUserParams) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return "", fmt.Errorf("[ERROR] %w: [%w]", pkg.ErrDbInternal, err)
	}

	if user.IsBanned.Bool {
		return "", pkg.ErrForbiden
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(payload.Password.String))
	if err != nil {
		return "", pkg.ErrWrongPassword
	}

	token, err := s.jwtHandler.GenUserToken(user.ID, user.UserType.UsersUserType, time.Now().Add(tokenTimeToLive))
	if err != nil {
		pkg.PrintErr(pkg.ErrCreatingToken, err)
		return "", pkg.ErrCreatingToken
	}

	return token, nil
}

func (s AuthService) LoginPet(ctx context.Context, petId int64, ownerId int64) (string, error) {
	pet, err := s.userRepo.GetPetById(ctx, petId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return "", fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	if pet.OwnerID.Int64 != ownerId {
		return "", pkg.ErrForbiden
	}

	token, err := s.jwtHandler.GenUserToken(pet.ID, core.UsersUserTypePet, time.Now().Add(time.Hour*24))
	if err != nil {
		pkg.PrintErr(pkg.ErrCreatingToken, err)
		return "", fmt.Errorf("%w: [%w]", pkg.ErrCreatingToken, err)
	}

	return token, nil
}

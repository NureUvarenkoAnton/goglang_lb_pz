package transport

import (
	core "NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/api"
	"NureUvarenkoAnton/unik_go_lb_4/internal/views/auth"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService iAuthService
}

type iAuthService interface {
	RegisterUser(ctx context.Context, user core.CreateUserParams, toHashPassword bool) (string, error)
	Login(ctx context.Context, payload core.CreateUserParams) (string, error)
	LoginPet(ctx context.Context, petId int64, ownerId int64) (string, error)
}

func NewAuthHandler(service iAuthService) *AuthHandler {
	return &AuthHandler{
		authService: service,
	}
}

func (h AuthHandler) RegisterUser(ctx *gin.Context) {
	var payload api.RegisterPayload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		fmt.Println("err: ", err)
		ctx.AbortWithError(http.StatusBadRequest, pkg.ErrPayloadDecode)
		return
	}

	if core.UsersUserType(payload.UserType) == core.UsersUserTypeAdmin {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	if payload.UserType == "" {
		payload.UserType = string(core.UsersUserTypeDefault)
	}

	token, err := h.authService.RegisterUser(ctx, core.CreateUserParams{
		Email:    sql.NullString{String: payload.Email, Valid: true},
		Name:     sql.NullString{String: payload.Name, Valid: true},
		Password: sql.NullString{String: payload.Password, Valid: true},
		UserType: core.NullUsersUserType{UsersUserType: core.UsersUserType(payload.UserType), Valid: true},
	}, true)
	if err != nil {
		if errors.Is(err, pkg.ErrEmailDuplicate) {
			ctx.AbortWithError(http.StatusConflict, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.SetCookie("authToken", token, 0, "", "localhost", true, true)

	ctx.JSON(http.StatusOK, nil)
}

func (h AuthHandler) RegisterForm(ctx *gin.Context) {
	auth.RegisterPage().Render(ctx, ctx.Writer)
}

func (h AuthHandler) LoginForm(ctx *gin.Context) {
	auth.LoginPage().Render(ctx, ctx.Writer)
}

func (h AuthHandler) Login(ctx *gin.Context) {
	type LoginPayload struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	var payload LoginPayload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrPayloadDecode)
		return
	}

	token, err := h.authService.Login(ctx, core.CreateUserParams{
		Email:    sql.NullString{String: payload.Email, Valid: true},
		Password: sql.NullString{String: payload.Password, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if errors.Is(err, pkg.ErrForbiden) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if errors.Is(err, pkg.ErrWrongPassword) {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Header("HX-Redirect", "/profile/")
	ctx.SetCookie("authToken", token, 0, "", "localhost", true, true)

	ctx.JSON(http.StatusOK, nil)
}

func (h AuthHandler) LoginPet(ctx *gin.Context) {
	ownerId := ctx.GetInt64("user_id")
	if ownerId == 0 {

		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type LoginPetPayload struct {
		PetId int64 `json:"petId"`
	}

	var payload LoginPetPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.authService.LoginPet(ctx, payload.PetId, ownerId)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if errors.Is(err, pkg.ErrForbiden) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, api.TokenResponse{
		Token: token,
	})
}

func (h AuthHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("authToken", "", 0, "", "localhost", true, true)
	ctx.Header("HX-Redirect", "/auth/login")
}

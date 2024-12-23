package transport

import (
	core "NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/api"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/file"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService   iUserService
	ratingService iRatingService
	authService   iAuthService
}

func NewUserHandler(userService iUserService, ratingService iRatingService, authService iAuthService) *UserHandler {
	return &UserHandler{
		userService:   userService,
		ratingService: ratingService,
		authService:   authService,
	}
}

type iUserService interface {
	MarkDeleted(ctx context.Context, id int64) error
	RestoreFromDeletion(ctx context.Context, id int64) error
	BanUser(ctx context.Context, id int64) error
	// returns all users if no parameters provided
	GetUsers(ctx context.Context, params core.GetUsersParams) ([]api.UserResponse, error)
	GetById(ctx context.Context, id int64, requesterType core.UsersUserType) (api.UserResponse, error)
}

func (h *UserHandler) GetUsersAdmin(ctx *gin.Context) {
	type GetUsersQueryParams struct {
		UserType  string `form:"userType,omitempty"`
		IsBanned  *bool  `form:"isBanned,omitempty"`
		IsDeleted *bool  `form:"isDeleted,omitempty"`
	}
	var payload GetUsersQueryParams
	err := ctx.ShouldBindQuery(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userType := core.UsersUserType(payload.UserType)

	users, err := h.userService.GetUsers(ctx, core.GetUsersParams{
		UserType:  core.NullUsersUserType{UsersUserType: userType, Valid: userType != ""},
		IsBanned:  sql.NullBool{Bool: payload.IsBanned != nil && *payload.IsBanned, Valid: payload.IsBanned != nil},
		IsDeleted: sql.NullBool{Bool: payload.IsDeleted != nil && *payload.IsDeleted, Valid: payload.IsDeleted != nil},
	})
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserById(ctx *gin.Context) {
	type IdUriParams struct {
		Id int64 `uri:"id"`
	}
	var payload IdUriParams
	err := ctx.ShouldBindUri(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userType := core.UsersUserType(ctx.GetString("user_type"))

	users, err := h.userService.GetById(ctx, payload.Id, userType)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if errors.Is(err, pkg.ErrForbiden) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetWalkers(ctx *gin.Context) {
	users, err := h.userService.GetUsers(ctx, core.GetUsersParams{
		UserType: core.NullUsersUserType{
			UsersUserType: core.UsersUserTypeWalker,
			Valid:         true,
		},
		IsBanned:  sql.NullBool{Bool: false, Valid: true},
		IsDeleted: sql.NullBool{Bool: false, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetSelf(ctx *gin.Context) {
	id := ctx.GetInt64("user_id")
	if id == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetById(ctx, id, core.UsersUserTypeAdmin)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	type DeleteUserPayload struct {
		Id int64 `uri:"id"`
	}
	var payload DeleteUserPayload
	err := ctx.ShouldBindUri(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.userService.MarkDeleted(ctx, payload.Id)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) DeleteSelf(ctx *gin.Context) {
	id := ctx.GetInt64("user_id")

	err := h.userService.MarkDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *UserHandler) RestoreFromDeletion(ctx *gin.Context) {
	id := ctx.GetInt64("user_id")

	err := h.userService.RestoreFromDeletion(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *UserHandler) SetBanState(ctx *gin.Context) {
	type BanUserPayload struct {
		Id int64 `json:"id"`
	}
	var payload BanUserPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.userService.BanUser(ctx, payload.Id)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) ExportUsers(ctx *gin.Context) {
	users, err := h.userService.GetUsers(ctx, core.GetUsersParams{})
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dataName := file.CreateFile(users, "json")
	ctx.FileAttachment(dataName, "users.json")
	ctx.Status(http.StatusOK)
	file.DeleteFile(dataName)
}

func (h *UserHandler) ImportUsers(ctx *gin.Context) {
	fmt.Println("retrieving file from formdata")
	header, err := ctx.FormFile("data")
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fmt.Println("2")
	fileDataReader, err := header.Open()
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fmt.Println("3")
	fileData, err := io.ReadAll(fileDataReader)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var users []api.UserResponse
	err = json.Unmarshal(fileData, &users)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return

	}

	for _, user := range users {
		h.authService.RegisterUser(ctx, core.CreateUserParams{
			Name:     sql.NullString{String: user.Name, Valid: true},
			Email:    sql.NullString{String: user.Email, Valid: true},
			UserType: core.NullUsersUserType{UsersUserType: core.UsersUserType(user.UserType), Valid: true},
		}, false)
	}

	ctx.Status(http.StatusOK)
}

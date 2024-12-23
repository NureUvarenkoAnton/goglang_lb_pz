package transport

import (
	core "NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/api"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/translate"
	"NureUvarenkoAnton/unik_go_lb_4/internal/views/profile"
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type ProfileHandler struct {
	profileService iUserProfileService
	userService    iUserService
	walksService   iWalkService
	m              *melody.Melody
}

func NewProfileHandler(
	service iUserProfileService,
	userService iUserService,
	walksService iWalkService,
	melody *melody.Melody,
) *ProfileHandler {
	return &ProfileHandler{
		profileService: service,
		userService:    userService,
		walksService:   walksService,
		m:              melody,
	}
}

type iUserProfileService interface {
	AddPet(ctx context.Context, pet core.AddPetParams) error
	GetPetById(ctx context.Context, id int64) (*core.Pet, error)
	UpdateUserData(ctx context.Context, userData core.UpdateUserParams) error
	UpdatePet(ctx context.Context, pet core.UpdatePetParams) error
	DeletePet(ctx context.Context, petId, ownerId int64) error
	GetAllPetsByOwnerId(ctx context.Context, lang string, ownerID sql.NullInt64) ([]core.Pet, error)
	GetTheMostWalkeblePetByOwnerID(ctx context.Context, ownerID sql.NullInt64) (core.Pet, error)
}

func (h *ProfileHandler) PetForm(ctx *gin.Context) {
	profile.PetForm().Render(ctx, ctx.Writer)
}

func (h *ProfileHandler) AddPet(ctx *gin.Context) {
	type AddPetPayload struct {
		Name           string `json:"name" form:"name"`
		Age            int    `json:"age" form:"age"`
		AdditionalInfo string `json:"additionalInfo" form:"additional_info"`
	}

	var payload AddPetPayload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	err = h.profileService.AddPet(ctx, core.AddPetParams{
		OwnerID:        sql.NullInt64{Int64: ctx.GetInt64("user_id"), Valid: true},
		Name:           sql.NullString{String: payload.Name, Valid: true},
		Age:            sql.NullInt16{Int16: int16(payload.Age), Valid: true},
		AdditionalInfo: sql.NullString{String: payload.AdditionalInfo, Valid: true},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.Header("HX-Location", "/profile")
}

func (h *ProfileHandler) GetOwnerPets(ctx *gin.Context) {
	type LangUri struct {
		Lang string `uri:"lang"`
	}
	var langPayload LangUri
	err := ctx.ShouldBindUri(&langPayload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if langPayload.Lang != translate.LANG_UA &&
		langPayload.Lang != translate.LANG_EN {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pets, err := h.profileService.GetAllPetsByOwnerId(ctx, langPayload.Lang, sql.NullInt64{
		Int64: ctx.GetInt64("user_id"),
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
	}

	ctx.JSON(http.StatusOK, api.SliceDbPetToApiPet(pets))
}

func (h *ProfileHandler) UpdatePet(ctx *gin.Context) {
	type UpdatePetPayload struct {
		PetId          int    `json:"pet_id"`
		Name           string `json:"name"`
		Age            int    `json:"age"`
		AdditionalInfo string `json:"additional_info"`
	}
	var payload UpdatePetPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	err = h.profileService.UpdatePet(ctx, core.UpdatePetParams{
		Name:           sql.NullString{String: payload.Name, Valid: payload.Name != ""},
		Age:            sql.NullInt16{Int16: int16(payload.Age), Valid: payload.Age != 0},
		AdditionalInfo: sql.NullString{String: payload.AdditionalInfo, Valid: payload.AdditionalInfo != ""},
		ID:             int64(payload.PetId),
	})
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
}

func (h *ProfileHandler) UpdateUser(ctx *gin.Context) {
	type UpdateUserPayload struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var payload UpdateUserPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	err = h.profileService.UpdateUserData(ctx, core.UpdateUserParams{
		Name:  sql.NullString{String: payload.Name, Valid: payload.Name != ""},
		Email: sql.NullString{String: payload.Email, Valid: payload.Email != ""},
		ID:    ctx.GetInt64("user_id"),
	})
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
}

func (h *ProfileHandler) DeltePet(ctx *gin.Context) {
	type IdUri struct {
		PetId int64 `uri:"id"`
	}
	var payload IdUri

	err := ctx.ShouldBindUri(&payload)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ownerId := ctx.GetInt64("user_id")
	err = h.profileService.DeletePet(ctx, payload.PetId, ownerId)
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

	ctx.Status(http.StatusOK)
}

package transport

import (
	core "NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"NureUvarenkoAnton/unik_go_lb_4/internal/views/profile"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ProfileHandler) UserProfilePage(ctx *gin.Context) {
	userType := ctx.GetString("user_type")

	switch userType {
	case string(core.UsersUserTypeDefault):
		h.defaultUserPage(ctx)
		return
	case string(core.UsersUserTypeWalker):
		h.walkerProfilePage(ctx)
		return
	}
}

func (h *ProfileHandler) defaultUserPage(ctx *gin.Context) {
	userId := ctx.GetInt64("user_id")
	userData, err := h.userService.GetById(ctx, userId, core.UsersUserTypeDefault)
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

	walks, err := h.walksService.GetWalksByOwnerId(ctx, sql.NullInt64{Int64: userId, Valid: true})
	if err != nil && !errors.Is(err, pkg.ErrNotFound) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	pets, err := h.profileService.GetAllPetsByOwnerId(ctx, "en", sql.NullInt64{Int64: userId, Valid: true})
	if err != nil && !errors.Is(err, pkg.ErrNotFound) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	theMostWalkeblePet, err := h.profileService.GetTheMostWalkeblePetByOwnerID(ctx, sql.NullInt64{Int64: userId, Valid: true})
	if err != nil && !errors.Is(err, pkg.ErrNotFound) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	result := profile.ProfildPageData{
		UserData: profile.UserData{
			Email: userData.Email,
			Name:  userData.Name,
			TheMostWalkeblePet: profile.IDContainer{
				ID:   fmt.Sprintf("%v", theMostWalkeblePet.ID),
				Name: theMostWalkeblePet.Name.String,
			},
		},
		Walks:        make([]profile.Walk, 0, len(walks)),
		PendingWalks: []profile.Walk{},
		Pets:         make([]profile.Pet, 0, len(pets)),
	}

	for _, walk := range walks {
		walkerData, err := h.userService.GetById(ctx, walk.WalkerID.Int64, core.UsersUserTypeDefault)
		if err != nil {
			continue
		}

		petData, err := h.profileService.GetPetById(ctx, walk.PetID.Int64)
		if err != nil {
			continue
		}

		result.Walks = append(result.Walks, profile.Walk{
			ID:         fmt.Sprintf("%d", walk.ID),
			WalkerName: walkerData.Name,
			PetName:    petData.Name.String,
		})
	}

	for _, pet := range pets {
		result.Pets = append(result.Pets, profile.Pet{
			ID:             fmt.Sprintf("%v", pet.ID),
			Name:           pet.Name.String,
			Age:            fmt.Sprintf("%v", pet.Age.Int16),
			AdditionalInfo: pet.AdditionalInfo.String,
		})
	}

	profile.ProfileDefaultPage(result).Render(ctx, ctx.Writer)
}

func (h *ProfileHandler) walkerProfilePage(ctx *gin.Context) {
	userId := ctx.GetInt64("user_id")
	userData, err := h.userService.GetById(ctx, userId, core.UsersUserTypeDefault)
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

	pendingWalks, err := h.walksService.GetWalksInfoByParams(ctx, core.GetWalkInfoByParamsParams{
		WalkerID:  sql.NullInt64{Int64: userId, Valid: true},
		WalkState: core.NullWalksState{WalksState: core.WalksStatePending, Valid: true},
	})
	if err != nil && !errors.Is(err, pkg.ErrNotFound) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	pendingWalksVM := make([]profile.Walk, 0, len(pendingWalks))

	for _, walk := range pendingWalks {
		pendingWalksVM = append(pendingWalksVM, profile.Walk{
			ID:        fmt.Sprintf("%v", walk.WalkID),
			OwnerName: fmt.Sprintf("%v", walk.OwnerName.String),
			PetName:   fmt.Sprintf("%v", walk.PetName.String),
		})
	}

	acceptedWalks, err := h.walksService.GetWalksInfoByParams(ctx, core.GetWalkInfoByParamsParams{
		WalkerID:  sql.NullInt64{Int64: userId, Valid: true},
		WalkState: core.NullWalksState{WalksState: core.WalksStateAccepted, Valid: true},
	})
	if err != nil && !errors.Is(err, pkg.ErrNotFound) {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	acceptedWalksVM := make([]profile.Walk, 0, len(acceptedWalks))

	for _, walk := range acceptedWalks {
		acceptedWalksVM = append(acceptedWalksVM, profile.Walk{
			ID:        fmt.Sprintf("%v", walk.WalkID),
			OwnerName: fmt.Sprintf("%v", walk.OwnerName.String),
			PetName:   fmt.Sprintf("%v", walk.PetName.String),
		})
	}

	profile.ProfileWalkerkPage(profile.ProfildPageData{
		UserData: profile.UserData{
			Email: userData.Email,
			Name:  userData.Name,
		},
		Walks:        acceptedWalksVM,
		PendingWalks: pendingWalksVM,
	}).Render(ctx, ctx.Writer)
}

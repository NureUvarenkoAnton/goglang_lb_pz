package transport

import (
	core "NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type PositionHandler struct {
	melody     *melody.Melody
	petService iPositionPetService
}

func NewPositionHandler(melody *melody.Melody, petService iPositionPetService) *PositionHandler {
	return &PositionHandler{
		melody:     melody,
		petService: petService,
	}
}

type iPositionPetService interface {
	GetPetById(ctx context.Context, id int64) (*core.Pet, error)
}

func (h *PositionHandler) HandleOpenPetConnection(ctx *gin.Context) {
	userId := ctx.GetInt64("user_id")
	userType := ctx.GetString("user_type")
	ownerId := int64(0)

	if core.UsersUserType(userType) == core.UsersUserTypePet {
		pet, err := h.petService.GetPetById(ctx, userId)
		if err != nil {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("pet doesn't exist"))
			return
		}
		ownerId = pet.OwnerID.Int64
	}

	// TODO: move to package
	h.melody.HandleRequestWithKeys(ctx.Writer, ctx.Request, map[string]any{
		"id":       userId,
		"userType": userType,
		"ownerId":  ownerId,
	})
}

func (h *PositionHandler) HandleMessage(s *melody.Session, b []byte) {
	userType := s.Keys["userType"].(string)

	switch core.UsersUserType(userType) {
	case core.UsersUserTypePet:
		ownerId := s.Keys["ownerId"]
		h.melody.BroadcastFilter(b, func(s *melody.Session) bool {
			fmt.Println(s.Keys["id"], " ", s.Keys["userType"])
			if s.Keys["id"] == ownerId && (s.Keys["userType"] == string(core.UsersUserTypeDefault)) {
				return true
			}
			return false
		})

	case core.UsersUserTypeDefault:
	case core.UsersUserTypeWalker:
	}
}

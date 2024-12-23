package middleware

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/jwt"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func TokenVerifier(jwtHandler jwt.JWT, userTypesAllowed []core.UsersUserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawToken, err := ctx.Cookie("authToken")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		if len(rawToken) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		claims, err := jwtHandler.VerifyToken(rawToken, userTypesAllowed)
		if err != nil {
			fmt.Println(err)
			if errors.Is(err, pkg.ErrForbiden) {
				ctx.AbortWithStatusJSON(http.StatusForbidden, nil)
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		if claims.ID == 0 ||
			!slices.Contains([]core.UsersUserType{
				core.UsersUserTypeAdmin,
				core.UsersUserTypeWalker,
				core.UsersUserTypeDefault,
				core.UsersUserTypePet,
			}, claims.UserType) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.Set("user_id", claims.ID)
		ctx.Set("user_type", string(claims.UserType))

		ctx.Next()
	}
}

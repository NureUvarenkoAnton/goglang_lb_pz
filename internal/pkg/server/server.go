package server

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/jwt"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/middleware"
	"NureUvarenkoAnton/unik_go_lb_4/internal/transport"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func New(handler *transport.Handler, jwtHandler jwt.JWT, melody *melody.Melody) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: setUpRoutes(handler, jwtHandler, melody),
	}
}

func setUpRoutes(handler *transport.Handler, jwtHandler jwt.JWT, melody *melody.Melody) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())

	authRouts := router.Group("/auth")
	{
		authRouts.GET("/register", handler.AuthHandler.RegisterForm)
		authRouts.POST("/register", handler.AuthHandler.RegisterUser)

		authRouts.GET("/login", handler.AuthHandler.LoginForm)
		authRouts.POST("/login", handler.AuthHandler.Login)

		authRouts.POST("/logout", handler.AuthHandler.Logout)
	}

	wsOpenConnection := router.Group("/")
	wsOpenConnection.Use(middleware.TokenVerifier(jwtHandler, []core.UsersUserType{core.UsersUserTypePet, core.UsersUserTypeDefault, core.UsersUserTypeWalker}))
	wsOpenConnection.GET("/ws", handler.PositionHandler.HandleOpenPetConnection)
	melody.HandleMessage(handler.PositionHandler.HandleMessage)

	loginPetRouter := router.Group("/")
	loginPetRouter.Use(middleware.TokenVerifier(jwtHandler, []core.UsersUserType{core.UsersUserTypeDefault}))
	loginPetRouter.POST("/loginpet", handler.AuthHandler.LoginPet)

	profileRouts := router.Group("/profile")
	profileRouts.Use(middleware.TokenVerifier(jwtHandler, []core.UsersUserType{
		core.UsersUserTypeAdmin,
		core.UsersUserTypeDefault,
		core.UsersUserTypeWalker,
	}))
	{
		profileRouts.GET("/", handler.ProfileHandler.UserProfilePage)

		profileRouts.GET("/pet/form", handler.ProfileHandler.PetForm)

		profileRouts.POST("/pet", handler.ProfileHandler.AddPet)
		profileRouts.PUT("/pet", handler.ProfileHandler.UpdatePet)
		profileRouts.GET("/pets/:lang", handler.ProfileHandler.GetOwnerPets)
		profileRouts.PUT("/user", handler.ProfileHandler.UpdateUser)
		profileRouts.DELETE("/pet/:id", handler.ProfileHandler.DeltePet)
	}

	usersDefaultRouts := router.Group("/users")
	usersDefaultRouts.Use(middleware.TokenVerifier(jwtHandler, []core.UsersUserType{
		core.UsersUserTypeAdmin,
		core.UsersUserTypeDefault,
		core.UsersUserTypeWalker,
	}))
	{
		usersDefaultRouts.GET("/self", handler.UserHandler.GetSelf)
		usersDefaultRouts.GET("/walkers", handler.UserHandler.GetWalkers)
		usersDefaultRouts.DELETE("/", handler.UserHandler.DeleteSelf)
		usersDefaultRouts.PUT("/restore", handler.UserHandler.RestoreFromDeletion)
	}

	userAdminRouts := router.Group("/users/admin")
	userAdminRouts.Use(middleware.TokenVerifier(jwtHandler, []core.UsersUserType{core.UsersUserTypeAdmin}))
	{
		userAdminRouts.GET("/:id", handler.UserHandler.GetUserById)
		userAdminRouts.PUT("/ban", handler.UserHandler.SetBanState)
		userAdminRouts.GET("/", handler.UserHandler.GetUsersAdmin)
		userAdminRouts.GET("/export", handler.UserHandler.ExportUsers)
		userAdminRouts.POST("/import", handler.UserHandler.ImportUsers)
	}

	walkRouts := router.Group("/walk")
	walkRouts.Use(middleware.TokenVerifier(jwtHandler, []core.UsersUserType{
		core.UsersUserTypeAdmin,
		core.UsersUserTypeWalker,
		core.UsersUserTypeDefault,
	}))
	{
		walkRouts.GET("/form", handler.WalkHalder.CreateWalkForm)
		walkRouts.POST("/", handler.WalkHalder.CreateWalkRequest)
		walkRouts.PUT("/", handler.WalkHalder.UpdateWalkState)
		walkRouts.GET("/:lang/:id", handler.WalkHalder.GetWalkInfoById)
		walkRouts.GET("/:lang", handler.WalkHalder.GetWalksByParams)
		walkRouts.GET("/self/:lang", handler.WalkHalder.GetWalksBySelfId)
	}

	ratingRouts := router.Group("/rating")
	ratingRouts.Use(middleware.TokenVerifier(jwtHandler, []core.UsersUserType{core.UsersUserTypeWalker, core.UsersUserTypeDefault, core.UsersUserTypeAdmin}))
	{
		ratingRouts.POST("/", handler.RatingHandler.AddRating)
		ratingRouts.GET("/:id", handler.RatingHandler.GetAvgRating)
	}

	return router
}

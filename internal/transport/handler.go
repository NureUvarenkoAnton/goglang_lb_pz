package transport

import "github.com/olahol/melody"

type Handler struct {
	AuthHandler     *AuthHandler
	ProfileHandler  *ProfileHandler
	PositionHandler *PositionHandler
	UserHandler     *UserHandler
	WalkHalder      *WalkHalder
	RatingHandler   *RatingHandler
}

func NewHandler(
	authService iAuthService,
	profileService iUserProfileService,
	melody *melody.Melody,
	userService iUserService,
	walkService iWalkService,
	ratingService iRatingService,
) *Handler {
	return &Handler{
		AuthHandler:     NewAuthHandler(authService),
		ProfileHandler:  NewProfileHandler(profileService, userService, walkService, melody),
		PositionHandler: NewPositionHandler(melody, profileService),
		UserHandler:     NewUserHandler(userService, ratingService, authService),
		WalkHalder:      NewWalkHandler(walkService, userService, profileService),
		RatingHandler:   NewRatingHandler(ratingService),
	}
}

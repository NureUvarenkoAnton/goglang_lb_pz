package api

import "NureUvarenkoAnton/unik_go_lb_4/internal/core"

type TokenResponse struct {
	Token string `json:"token"`
}

type ResponseWSMessage struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Data  any    `json:"data"`
}

type AvgRatingResponse struct {
	AvgRating float64 `json:"avgRating"`
}

type UserResponse struct {
	Id        int64              `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	UserType  core.UsersUserType `json:"userType"`
	AvgRating float64            `json:"avgRating,omitempty"`
	IsBanned  bool               `json:"isBanned"`
	IsDeleted bool               `json:"isDelted"`
}

type WalkInfoResponse struct {
	WalkId            int64  `json:"walkId,omitempty"`
	WalkState         string `json:"walkState,omitempty"`
	StartTime         string `json:"startTime,omitempty"`
	FinishTime        string `json:"finishTime,omitempty"`
	OwnerId           int64  `json:"ownerId,omitempty"`
	OwnerName         string `json:"ownerName,omitempty"`
	OwnerEmail        string `json:"ownerEmail,omitempty"`
	WalkerId          int64  `json:"walkerId,omitempty"`
	WalkerName        string `json:"walkerName,omitempty"`
	WalkerEmail       string `json:"walkerEmail,omitempty"`
	PetId             int64  `json:"petId,omitempty"`
	PetName           string `json:"petName,omitempty"`
	PetAdditionalInfo string `json:"PetAdditionalInfo,omitempty"`
}

type PetResponse struct {
	PetId             int64  `json:"petId,omitempty"`
	OwnerId           int64  `json:"ownerId"`
	Age               int64  `json:"age"`
	PetName           string `json:"petName,omitempty"`
	PetAdditionalInfo string `json:"PetAdditionalInfo,omitempty"`
}

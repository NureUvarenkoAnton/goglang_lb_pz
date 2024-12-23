package api

import "NureUvarenkoAnton/unik_go_lb_4/internal/core"

func DbUserToAPIUser(user core.User) UserResponse {
	return UserResponse{
		Id:        user.ID,
		Name:      user.Name.String,
		Email:     user.Email.String,
		UserType:  user.UserType.UsersUserType,
		IsBanned:  user.IsBanned.Bool,
		IsDeleted: user.IsDeleted.Bool,
	}
}

func SliceDbUserToAPIUser(users []core.User) []UserResponse {
	var result []UserResponse
	for _, user := range users {
		result = append(result, DbUserToAPIUser(user))
	}
	return result
}

func DbWalkInfoToAPIWalkInfo(walkInfo core.WalkInfo) WalkInfoResponse {
	result := WalkInfoResponse{
		WalkId:            walkInfo.WalkID,
		OwnerId:           walkInfo.OwnerID,
		OwnerName:         walkInfo.OwnerName.String,
		OwnerEmail:        walkInfo.OwnerEmail.String,
		WalkerId:          walkInfo.WalkerID,
		WalkerName:        walkInfo.WalkerName.String,
		WalkerEmail:       walkInfo.WalkerEmail.String,
		PetId:             walkInfo.PetID,
		PetName:           walkInfo.PetName.String,
		PetAdditionalInfo: walkInfo.PetAdditionalInfo.String,
		WalkState:         string(walkInfo.State.WalksState),
	}

	if walkInfo.StartTime.Valid {
		result.StartTime = walkInfo.StartTime.Time.String()
	}

	if walkInfo.FinishTime.Valid {
		result.FinishTime = walkInfo.FinishTime.Time.String()
	}

	return result
}

func SliceDbWalkInfoToAPIWalkInfo(walksInfo []core.WalkInfo) []WalkInfoResponse {
	var result []WalkInfoResponse
	for _, walkInfo := range walksInfo {
		result = append(result, DbWalkInfoToAPIWalkInfo(walkInfo))
	}
	return result
}

func DbPetToApiPet(pet core.Pet) PetResponse {
	return PetResponse{
		PetId:             pet.ID,
		OwnerId:           pet.OwnerID.Int64,
		Age:               int64(pet.Age.Int16),
		PetName:           pet.Name.String,
		PetAdditionalInfo: pet.AdditionalInfo.String,
	}
}

func SliceDbPetToApiPet(pets []core.Pet) []PetResponse {
	var result []PetResponse
	for _, pet := range pets {
		result = append(result, DbPetToApiPet(pet))
	}
	return result
}

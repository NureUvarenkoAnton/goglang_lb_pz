package service

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type ProfileService struct {
	userRepo      iProfileRepo
	petTranslator iPetTranlsator
}

func NewProfileService(userRepo iProfileRepo, petTranslator iPetTranlsator) *ProfileService {
	return &ProfileService{
		userRepo:      userRepo,
		petTranslator: petTranslator,
	}
}

type iProfileRepo interface {
	GetPetById(ctx context.Context, id int64) (core.Pet, error)
	UpdatePet(ctx context.Context, arg core.UpdatePetParams) error
	GetAllPetsByOwnerId(ctx context.Context, ownerID sql.NullInt64) ([]core.Pet, error)
	AddPet(ctx context.Context, arg core.AddPetParams) error
	DeletePet(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, arg core.UpdateUserParams) error
	GetTheMostWalkeblePetByOwnerID(ctx context.Context, ownerID int64) (core.GetTheMostWalkeblePetByOwnerIDRow, error)
}

type iPetTranlsator interface {
	Tranlate(text, targetLang string) string
}

func (s *ProfileService) AddPet(ctx context.Context, pet core.AddPetParams) error {
	err := s.userRepo.AddPet(ctx, pet)
	if err != nil {
		err, ok := err.(*mysql.MySQLError)
		if !ok {
			pkg.PrintErr(pkg.ErrDbInternal, err)
			return pkg.ErrDbInternal
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return pkg.ErrDbInternal
	}

	return err
}

func (s *ProfileService) GetPetById(ctx context.Context, id int64) (*core.Pet, error) {
	pet, err := s.userRepo.GetPetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return nil, pkg.ErrDbInternal
	}

	return &pet, nil
}

func (s *ProfileService) GetAllPetsByOwnerId(ctx context.Context, lang string, ownerID sql.NullInt64) ([]core.Pet, error) {
	pets, err := s.userRepo.GetAllPetsByOwnerId(ctx, ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return nil, pkg.ErrDbInternal
	}

	return pets, nil
}

func (s *ProfileService) GetTheMostWalkeblePetByOwnerID(ctx context.Context, ownerID sql.NullInt64) (core.Pet, error) {
	raw, err := s.userRepo.GetTheMostWalkeblePetByOwnerID(ctx, ownerID.Int64)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Pet{}, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return core.Pet{}, pkg.ErrDbInternal
	}

	return core.Pet{
		ID:      raw.PetID,
		OwnerID: ownerID,
		Name:    sql.NullString{String: raw.PetName.String, Valid: true},
	}, nil
}

func (s *ProfileService) UpdatePet(ctx context.Context, pet core.UpdatePetParams) error {
	_, err := s.userRepo.GetPetById(ctx, pet.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return err
	}

	err = s.userRepo.UpdatePet(ctx, pet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return pkg.ErrDbInternal

	}

	return nil
}

func (s *ProfileService) DeletePet(ctx context.Context, petId, ownerId int64) error {
	pet, err := s.userRepo.GetPetById(ctx, petId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	if pet.OwnerID.Int64 != ownerId {
		return pkg.ErrForbiden
	}

	err = s.userRepo.DeletePet(ctx, petId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return nil
}

func (s *ProfileService) UpdateUserData(ctx context.Context, userData core.UpdateUserParams) error {
	err := s.userRepo.UpdateUser(ctx, userData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		err, ok := err.(*mysql.MySQLError)
		if !ok {
			pkg.PrintErr(pkg.ErrDbInternal, err)
			return pkg.ErrDbInternal
		}
		if err.Number == 1062 {
			return pkg.ErrEmailDuplicate
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return pkg.ErrDbInternal
	}
	return nil
}

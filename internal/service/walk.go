package service

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type WalkService struct {
	walkRepo       iWalkRepo
	walkTranslator iWalkTranslator
}

type iWalkTranslator interface {
	Tranlate(text, targetLang string) string
}

func NewWalkService(walkRepo iWalkRepo, walkTranslator iWalkTranslator) *WalkService {
	return &WalkService{
		walkRepo:       walkRepo,
		walkTranslator: walkTranslator,
	}
}

type iWalkRepo interface {
	GetPetById(ctx context.Context, id int64) (core.Pet, error)
	GetUserById(ctx context.Context, id int64) (core.User, error)
	CreateWalk(ctx context.Context, arg core.CreateWalkParams) error
	GetWalksByWalkerId(ctx context.Context, walkerID sql.NullInt64) ([]core.Walk, error)
	UpdateWalkState(ctx context.Context, arg core.UpdateWalkStateParams) error
	GetWalksByOwnerId(ctx context.Context, ownerID sql.NullInt64) ([]core.Walk, error)
	GetWalkById(ctx context.Context, id int64) (core.Walk, error)
	GetWalkInfoByParams(ctx context.Context, arg core.GetWalkInfoByParamsParams) ([]core.WalkInfo, error)
	GetWalkInfoByWalkId(ctx context.Context, walkID int64) (core.WalkInfo, error)
}

func (s *WalkService) CreateWalk(ctx context.Context, walkParams core.CreateWalkParams) error {
	walker, err := s.walkRepo.GetUserById(ctx, walkParams.WalkerID.Int64)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	if walker.UserType.UsersUserType != core.UsersUserTypeWalker {
		return pkg.ErrNotFound
	}

	pet, err := s.walkRepo.GetPetById(ctx, walkParams.PetID.Int64)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	if pet.OwnerID.Int64 != walkParams.OwnerID.Int64 {
		return pkg.ErrForbiden
	}

	err = s.walkRepo.CreateWalk(ctx, walkParams)
	if err != nil {
		err, ok := err.(*mysql.MySQLError)
		if !ok {
			pkg.PrintErr(pkg.ErrDbInternal, err)
			return pkg.ErrDbInternal
		}

		return pkg.ErrDbInternal
	}
	return nil
}

func (s *WalkService) GetWalksByWalkerId(ctx context.Context, walkerID sql.NullInt64) ([]core.Walk, error) {
	walks, err := s.walkRepo.GetWalksByWalkerId(ctx, walkerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return nil, pkg.ErrDbInternal
	}
	return walks, nil
}

func (s *WalkService) GetWalksByOwnerId(ctx context.Context, ownerID sql.NullInt64) ([]core.Walk, error) {
	walks, err := s.walkRepo.GetWalksByOwnerId(ctx, ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return nil, pkg.ErrDbInternal
	}

	return walks, nil
}

func (s *WalkService) UpdateWalkState(ctx context.Context, params core.UpdateWalkStateParams) error {
	_, err := s.walkRepo.GetWalkById(ctx, params.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	if params.State.WalksState == core.WalksStateFinished {
		params.FinishTime.Time = time.Now()
		params.FinishTime.Valid = true
	}

	err = s.walkRepo.UpdateWalkState(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return pkg.ErrDbInternal
	}
	return nil
}

func (s *WalkService) GetWalksInfoByParams(
	ctx context.Context,
	params core.GetWalkInfoByParamsParams,
) ([]core.WalkInfo, error) {
	info, err := s.walkRepo.GetWalkInfoByParams(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return nil, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return info, nil
}

func (s *WalkService) GetWalkInfoByWalkId(ctx context.Context, lang string, walkId int64) (core.WalkInfo, error) {
	walkInfo, err := s.walkRepo.GetWalkInfoByWalkId(ctx, walkId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.WalkInfo{}, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return core.WalkInfo{}, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return walkInfo, nil
}

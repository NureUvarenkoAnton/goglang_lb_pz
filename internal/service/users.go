package service

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/api"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/statistics"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserService struct {
	userRepo iUsersRepo
}

func NewUserService(userRepo iUsersRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type iUsersRepo interface {
	GetAllUsers(ctx context.Context) ([]core.User, error)
	SetBanState(ctx context.Context, arg core.SetBanStateParams) error
	SetDeleteState(ctx context.Context, arg core.SetDeleteStateParams) error
	DeleteMarkedUsers(ctx context.Context, deletedAt sql.NullTime) error
	GetUserById(ctx context.Context, id int64) (core.User, error)
	GetUsers(ctx context.Context, arg core.GetUsersParams) ([]core.User, error)
	RatingsByRateeId(ctx context.Context, rateeID sql.NullInt64) ([]core.Rating, error)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]api.UserResponse, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return nil, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	var usersResponse []api.UserResponse
	for _, user := range users {
		ratings, err := s.userRepo.RatingsByRateeId(ctx, sql.NullInt64{Int64: int64(user.ID), Valid: true})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
		}
		var ratingValues []int32
		for _, rating := range ratings {
			ratingValues = append(ratingValues, rating.Value.Int32)
		}

		usersResponse = append(usersResponse, api.UserResponse{
			Id:        user.ID,
			Name:      user.Name.String,
			Email:     user.Email.String,
			UserType:  user.UserType.UsersUserType,
			AvgRating: statistics.AvgWeighted(ratingValues),
			IsBanned:  user.IsBanned.Bool,
			IsDeleted: user.IsBanned.Bool,
		})
	}

	return usersResponse, nil
}

func (s *UserService) GetById(ctx context.Context, id int64, requesterType core.UsersUserType) (api.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.UserResponse{}, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return api.UserResponse{}, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	ratings, err := s.userRepo.RatingsByRateeId(ctx, sql.NullInt64{Int64: int64(user.ID), Valid: true})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		pkg.PrintErr(pkg.ErrDbInternal, err)
		return api.UserResponse{}, fmt.Errorf("%w: %w", pkg.ErrDbInternal, err)
	}
	var ratingValues []int32
	for _, rating := range ratings {
		ratingValues = append(ratingValues, rating.Value.Int32)
	}

	return api.UserResponse{
		Id:        user.ID,
		Name:      user.Name.String,
		Email:     user.Email.String,
		UserType:  user.UserType.UsersUserType,
		AvgRating: statistics.AvgWeighted(ratingValues),
		IsBanned:  user.IsBanned.Bool,
		IsDeleted: user.IsBanned.Bool,
	}, nil
}

func (s *UserService) GetUsers(ctx context.Context, params core.GetUsersParams) ([]api.UserResponse, error) {
	// if no paramters provided, then return all users
	if !params.IsBanned.Valid &&
		!params.UserType.Valid &&
		!params.IsDeleted.Valid {

		return s.GetAllUsers(ctx)
	}

	users, err := s.userRepo.GetUsers(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return nil, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	var usersResponse []api.UserResponse
	for _, user := range users {
		ratings, err := s.userRepo.RatingsByRateeId(ctx, sql.NullInt64{Int64: int64(user.ID), Valid: true})
		if err != nil {
			return nil, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
		}
		var ratingValues []int32
		for _, rating := range ratings {
			ratingValues = append(ratingValues, rating.Value.Int32)
		}

		usersResponse = append(usersResponse, api.UserResponse{
			Id:        user.ID,
			Name:      user.Name.String,
			Email:     user.Email.String,
			UserType:  user.UserType.UsersUserType,
			AvgRating: statistics.AvgWeighted(ratingValues),
			IsBanned:  user.IsBanned.Bool,
			IsDeleted: user.IsBanned.Bool,
		})
	}

	return usersResponse, nil
}

func (s *UserService) BanUser(ctx context.Context, id int64) error {
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	err = s.userRepo.SetBanState(ctx, core.SetBanStateParams{
		IsBanned: sql.NullBool{Bool: !user.IsBanned.Bool, Valid: true},
		ID:       user.ID,
	})
	if err != nil {
		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return nil
}

func (s *UserService) MarkDeleted(ctx context.Context, id int64) error {
	err := s.userRepo.SetDeleteState(ctx, core.SetDeleteStateParams{
		ID:        id,
		IsDeleted: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}
	return nil
}

func (s *UserService) DelteMarkedUsers(ctx context.Context, t time.Time) {
	err := s.userRepo.DeleteMarkedUsers(ctx, sql.NullTime{Time: t, Valid: true})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		pkg.PrintErr(pkg.ErrDbInternal, err)
	}
}

func (s *UserService) RestoreFromDeletion(ctx context.Context, id int64) error {
	err := s.userRepo.SetDeleteState(ctx, core.SetDeleteStateParams{
		ID:        id,
		IsDeleted: sql.NullBool{Bool: false, Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNotFound
		}

		pkg.PrintErr(pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}
	return nil
}

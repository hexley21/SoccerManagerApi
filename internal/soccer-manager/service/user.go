package service

import (
	"context"
	"errors"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	"github.com/hexley21/soccer-manager/pkg/hasher"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=mock/mock_user.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/service UserService
type UserService interface {
	Get(ctx context.Context, userId int64) (domain.User, error)
	List(ctx context.Context, cursor int64, limit int32) ([]domain.User, error)
	UpdatePassword(ctx context.Context, id int64, oldPassowrd string, newPassword string) error
	Delete(ctx context.Context, userId int64) error
}

type userServiceImpl struct {
	userRepo repository.UserRepository
	hasher   hasher.Hasher
}

func NewUserService(userRepo repository.UserRepository, hasher hasher.Hasher) *userServiceImpl {
	return &userServiceImpl{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

// Get finds a single user by id
//
// If not found - ErrUserNotFound
func (s *userServiceImpl) Get(ctx context.Context, userId int64) (domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}

		return domain.User{}, err
	}

	return domain.NewUser(userId, user.Username, user.Role), nil
}

// List returns a slice of users by pagination parameters
//
// Always returns empty slice
func (s *userServiceImpl) List(
	ctx context.Context,
	cursor int64,
	limit int32,
) ([]domain.User, error) {
	users, err := s.userRepo.ListUsersCursor(ctx, repository.ListUsersCursorParams{
		ID:    cursor,
		Limit: limit,
	})
	if err != nil {
		return []domain.User{}, err
	}

	res := make([]domain.User, len(users))

	for i, usr := range users {
		res[i] = domain.NewUser(usr.ID, usr.Username, usr.Role)
	}

	return res, nil
}

// UpdatePassword updates password by validating old one first and updating after
//
// If user not found - ErrUserNotFound
// If incorrect password - ErrIncorrectPassword
func (s *userServiceImpl) UpdatePassword(
	ctx context.Context,
	id int64,
	oldPassowrd string,
	newPassword string,
) error {
	oldHash, err := s.userRepo.GetUserHashByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	err = s.hasher.VerifyPassword(oldPassowrd, oldHash)
	if err != nil {
		return err
	}

	newHash, err := s.hasher.HashPassword(newPassword)
	if err != nil {
		if errors.Is(err, hasher.ErrPasswordMismatch) {
			return ErrIncorrectPassword
		}
		return err
	}

	err = s.userRepo.UpdateUserHash(ctx, repository.UpdateUserHashParams{
		ID:   id,
		Hash: newHash,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	return nil
}

// Delete removes user form db
//
// If user not found - ErrUserNotFound
func (s *userServiceImpl) Delete(ctx context.Context, userId int64) error {
	err := s.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	return nil
}

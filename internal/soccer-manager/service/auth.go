package service

import (
	"context"
	"errors"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	"github.com/hexley21/soccer-manager/pkg/hasher"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthService interface {
	Authenticate(ctx context.Context, username string, password string) (domain.User, error)
	CreateUser(
		ctx context.Context,
		username string,
		password string,
		role string,
	) (domain.User, error)
	GetUserById(ctx context.Context, userID int64) (domain.User, error)
}

type authServiceImpl struct {
	userRepo repository.UserRepository
	hasher   hasher.Hasher
}

func NewAuthService(userRepo repository.UserRepository, hasher hasher.Hasher) *authServiceImpl {
	return &authServiceImpl{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

// Authenticate user by searching in db and verifying password
//
// If user was not found: ErrUserNotfound
// If password was incorrect: ErrIncorrectPassword
func (s *authServiceImpl) Authenticate(
	ctx context.Context,
	username string,
	password string,
) (domain.User, error) {
	auth, err := s.userRepo.GetAuth(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}

		return domain.User{}, err
	}

	if err := s.hasher.VerifyPassword(password, auth.Hash); err != nil {
		return domain.User{}, ErrIncorrectPassword
	}

	return domain.NewUser(auth.ID, username, auth.Role), nil
}

// CreateUser by generating password hash & inserting into db
//
// If username is taken: ErrUsernameTaken
func (s *authServiceImpl) CreateUser(
	ctx context.Context,
	username string,
	password string,
	role string,
) (domain.User, error) {
	hash, err := s.hasher.HashPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	user, err := s.userRepo.CreateUser(ctx, repository.CreateUserParams{
		Username: username,
		Role:     role,
		Hash:     hash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.User{}, ErrUsernameTaken
		}

		return domain.User{}, err
	}

	return domain.NewUser(user.ID, user.Username, user.Role), nil
}

// GetUserById gets a single user from db
//
// If not found - ErrUserNotFound
func (s *authServiceImpl) GetUserById(ctx context.Context, userID int64) (domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}

		return domain.User{}, err
	}

	return domain.NewUser(user.ID, user.Username, user.Role), nil
}

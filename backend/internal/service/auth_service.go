package service

import (
	"database/sql"
	"errors"

	"github.com/itlabil/smmusic/backend/internal/config"
	"github.com/itlabil/smmusic/backend/internal/models"
	"github.com/itlabil/smmusic/backend/internal/repository"
	"github.com/itlabil/smmusic/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

// Login validates credentials and returns a signed JWT
func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID, user.Role, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// CreateUser is admin-only: creates a new user with a hashed password
func (s *AuthService) CreateUser(username, password, role string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword lets a logged-in user update their own password
func (s *AuthService) ChangePassword(userID int, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return sql.ErrNoRows
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, string(newHash))
}

func (s *AuthService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}

// AdminUpdateUser lets an admin change another user's role and/or reset their password.
// Both newPassword and newRole are optional (empty string = no change).
func (s *AuthService) AdminUpdateUser(userID int, newPassword, newRole string) error {
	if newRole != "" {
		if err := s.userRepo.UpdateRole(userID, newRole); err != nil {
			return err
		}
	}

	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := s.userRepo.UpdatePassword(userID, string(hash)); err != nil {
			return err
		}
	}

	return nil
}
package usecase

import (
	"errors"
	"log"

	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(username, password, role string) error
	Login(username, password string) (*domain.User, error)
}

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Register(username, password, role string) error {
	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashed),
		Role:         role,
	}
	return u.repo.Create(user)
}

func (u *userUsecase) Login(username, password string) (*domain.User, error) {
	user, err := u.repo.FindByUsername(username)
	log.Println("Login query result:", user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

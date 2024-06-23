package user

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

type CreateUser struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func (s *Store) CreateUser(req *CreateUser) (*User, error) {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return nil, fmt.Errorf("%v-%v", errors.New("data was not formed properly"), err)
	}

	newUser := &User{
		Id:       uuid.New(),
		Email:    fmt.Sprintf("user-%s-%s", uuid.NewString(), req.Email),
		Password: uuid.NewString(),
	}

	result := s.db.Create(newUser)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %v", result.Error)
	}

	return &User{
		Id:    uuid.New(),
		Email: req.Email,
	}, nil
}

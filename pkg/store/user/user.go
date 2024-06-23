package user

import (
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email    string    `json:"email" gorm:"type:varchar(255);unique;not null" validate:"required,email"`
	Password string    `json:"password" gorm:"type:varchar(255);not null" validate:"required,min=8"`
}

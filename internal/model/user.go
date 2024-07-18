package model

import "github.com/google/uuid"

type User struct {
	CommonBase
	Email          string   `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password       string   `json:"password" gorm:"type:varchar(255)"`
	OAuthProviders []string `json:"oauth_providers" gorm:"type:text[]"`
	IsActive       bool     `json:"is_active" gorm:"default:false"`
}

type Profile struct {
	CommonBase
	UserId    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	FullName  string    `json:"full_name" gorm:"type:varchar(255)"`
	FirstName string    `json:"first_name" gorm:"type:varchar(255)"`
	LastName  string    `json:"last_name" gorm:"type:varchar(255)"`
	AvatarURL string    `json:"avatar_url" gorm:"type:text"`
}

type AuthResponse struct {
	CommonBase
	UserId       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;"`
	AccessToken  string    `json:"access_token" gorm:"type:text;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"type:text;not null"`
	ExpiresIn    int64     `json:"expires_in" gorm:"type:bigint;not null"`
	TokenType    string    `json:"token_type" gorm:"type:varchar(50);not null;default:'Bearer'"`
}

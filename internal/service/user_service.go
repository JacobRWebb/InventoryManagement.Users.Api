package service

import (
	"context"
	"errors"
	"time"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/model"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/util"
	UserProto "github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/api"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ERR_INTERNAL_TRYAGAIN   = "Internal server error try again, please try again!"
	ERR_INVALID_CREDENTIALS = "Invalid credentials"
	ERR_USER_NOT_FOUND      = "Account was not found"
	ERR_EMAIL_TAKEN         = "Email is already registered"
	ERR_INVALID_TOKEN       = "Invalid token"
)

type UserService struct {
	db *gorm.DB
	UserProto.UnimplementedUserServiceServer
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req *UserProto.RegisterUserRequest) (*UserProto.Empty, error) {
	var existingUser model.User
	if err := s.db.Where("email = ?", req.GetEmail()).First(&existingUser).Error; err == nil {
		return nil, status.Error(codes.AlreadyExists, ERR_EMAIL_TAKEN)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	newUser := &model.User{
		Email:    req.GetEmail(),
		Password: string(hashedPassword),
		IsActive: true,
	}

	if err := tx.Create(newUser).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	newProfile := &model.Profile{
		UserId: newUser.Id,
	}

	if err := tx.Create(newProfile).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	authResponse, err := util.CreateAuthResponse(newUser.Id)

	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	if err := tx.Create(authResponse).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	return nil, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *UserProto.LoginUserRequest) (*UserProto.AuthResponse, error) {
	var user model.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ERR_USER_NOT_FOUND)
		}
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, ERR_INVALID_CREDENTIALS)
	}

	authResponse, err := util.CreateAuthResponse(user.Id)

	if err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	if err := s.db.Create(authResponse).Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	return &UserProto.AuthResponse{
		AccessToken:  authResponse.AccessToken,
		RefreshToken: authResponse.RefreshToken,
		ExpiresIn:    authResponse.ExpiresIn,
		TokenType:    authResponse.TokenType,
	}, nil
}

func (s *UserService) LoginWithOAuth(ctx context.Context, req *UserProto.OAuthLoginRequest) (*UserProto.AuthResponse, error) {
	// Implement OAuth login logic here
	// This would typically involve verifying the OAuth token with the provider
	// and either creating a new user or logging in an existing user
	return nil, status.Error(codes.Unimplemented, "OAuth login not implemented")
}

func (s *UserService) LogoutUser(ctx context.Context, req *UserProto.LogoutRequest) (*UserProto.LogoutResponse, error) {
	if err := s.db.Where("access_token = ?", req.AccessToken).Delete(&model.AuthResponse{}).Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}
	return &UserProto.LogoutResponse{Success: true}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req *UserProto.RefreshTokenRequest) (*UserProto.AuthResponse, error) {
	var authResponse model.AuthResponse
	if err := s.db.Where("refresh_token = ?", req.RefreshToken).First(&authResponse).Error; err != nil {
		return nil, status.Error(codes.Unauthenticated, ERR_INVALID_TOKEN)
	}

	newAuthResponse, err := util.CreateAuthResponse(authResponse.UserId)

	if err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	if err := s.db.Save(&newAuthResponse).Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	return &UserProto.AuthResponse{
		AccessToken:  newAuthResponse.AccessToken,
		RefreshToken: newAuthResponse.RefreshToken,
		ExpiresIn:    newAuthResponse.ExpiresIn,
		TokenType:    newAuthResponse.TokenType,
	}, nil
}

func (s *UserService) RevokeToken(ctx context.Context, req *UserProto.RevokeTokenRequest) (*UserProto.RevokeTokenResponse, error) {
	var authResponse model.AuthResponse
	query := s.db.Where("access_token = ?", req.Token)
	if req.TokenTypeHint == UserProto.TokenType_REFRESH_TOKEN {
		query = s.db.Where("refresh_token = ?", req.Token)
	}

	if err := query.Delete(&authResponse).Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	return &UserProto.RevokeTokenResponse{Success: true}, nil
}

func (s *UserService) ValidateToken(ctx context.Context, req *UserProto.ValidateTokenRequest) (*UserProto.ValidateTokenResponse, error) {
	_, claims, err := util.VerifyToken(req.GetAccessToken())

	if err != nil {
		return nil, err
	}

	userId, err := util.GetUserIdFromToken(claims)

	if err != nil {
		return nil, err
	}

	return &UserProto.ValidateTokenResponse{
		UserId:  userId.String(),
		IsValid: true,
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *UserProto.GetUserProfileRequest) (*UserProto.Profile, error) {
	var profile model.Profile
	if err := s.db.Where("user_id = ?", req.UserId).First(&profile).Error; err != nil {
		return nil, status.Error(codes.NotFound, ERR_USER_NOT_FOUND)
	}

	return &UserProto.Profile{
		UserId:    profile.UserId.String(),
		FullName:  profile.FullName,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		AvatarUrl: profile.AvatarURL,
	}, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, req *UserProto.UpdateUserProfileRequest) (*UserProto.Profile, error) {
	var profile model.Profile
	if err := s.db.Where("user_id = ?", req.UserId).First(&profile).Error; err != nil {
		return nil, status.Error(codes.NotFound, ERR_USER_NOT_FOUND)
	}

	profile.FullName = req.Profile.FullName
	profile.FirstName = req.Profile.FirstName
	profile.LastName = req.Profile.LastName
	profile.AvatarURL = req.Profile.AvatarUrl

	if err := s.db.Save(&profile).Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	return &UserProto.Profile{
		UserId:    profile.UserId.String(),
		FullName:  profile.FullName,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		AvatarUrl: profile.AvatarURL,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *UserProto.GetUserRequest) (*UserProto.User, error) {
	var user model.User
	if err := s.db.Where("id = ?", req.UserId).First(&user).Error; err != nil {
		return nil, status.Error(codes.NotFound, ERR_USER_NOT_FOUND)
	}

	var profile model.Profile
	if err := s.db.Where("user_id = ?", user.Id).First(&profile).Error; err != nil {
		return nil, status.Error(codes.NotFound, ERR_USER_NOT_FOUND)
	}

	return &UserProto.User{
		Id:    user.Id.String(),
		Email: user.Email,
		Profile: &UserProto.Profile{
			UserId:    profile.UserId.String(),
			FullName:  profile.FullName,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			AvatarUrl: profile.AvatarURL,
		},
		OauthProviders: user.OAuthProviders,
		IsActive:       user.IsActive,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *UserProto.ListUsersRequest) (*UserProto.ListUserResponse, error) {
	var users []model.User
	var totalCount int64

	offset := (req.Page - 1) * req.PageSize

	if err := s.db.Model(&model.User{}).Count(&totalCount).Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	if err := s.db.Offset(int(offset)).Limit(int(req.PageSize)).Find(&users).Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	var protoUsers []*UserProto.User
	for _, user := range users {
		var profile model.Profile
		if err := s.db.Where("user_id = ?", user.Id).First(&profile).Error; err != nil {
			return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
		}

		protoUsers = append(protoUsers, &UserProto.User{
			Id:    user.Id.String(),
			Email: user.Email,
			Profile: &UserProto.Profile{
				UserId:    profile.UserId.String(),
				FullName:  profile.FullName,
				FirstName: profile.FirstName,
				LastName:  profile.LastName,
				AvatarUrl: profile.AvatarURL,
			},
			OauthProviders: user.OAuthProviders,
			IsActive:       user.IsActive,
			CreatedAt:      user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      user.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &UserProto.ListUserResponse{
		Users:      protoUsers,
		TotalCount: int32(totalCount),
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *UserProto.DeleteUserRequest) (*UserProto.DeleteUserResponse, error) {
	if err := s.db.Where("id = ?", req.UserId).Delete(&model.User{}).Error; err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	return &UserProto.DeleteUserResponse{Success: true}, nil
}

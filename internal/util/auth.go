package util

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/JacobRWebb/InventoryManagement.Users.Api/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ERR_INTERNAL_TRYAGAIN   = "Internal server error try again, please try again!"
	ERR_INVALID_CREDENTIALS = "Invalid credentials"
	ERR_USER_NOT_FOUND      = "Account was not found"
	ERR_EMAIL_TAKEN         = "Email is already registered"
	ERR_INVALID_TOKEN       = "Invalid token"
)

func CreateAuthResponse(userId uuid.UUID) (*model.AuthResponse, error) {
	accessToken, refreshToken, err := generateTokens(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, ERR_INTERNAL_TRYAGAIN)
	}

	authResponse := &model.AuthResponse{
		UserId:       userId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}

	return authResponse, nil
}

func VerifyToken(tokenString string) (*jwt.Token, *jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your-256-bit-secret"), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, status.Error(codes.InvalidArgument, "token is not valid")
	}

	if exp, ok := (*claims)["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, nil, status.Error(codes.InvalidArgument, "token is not valid")
		}
	}

	return token, claims, nil
}

func GetUserIdFromToken(claims *jwt.MapClaims) (uuid.UUID, error) {
	userIdStr, ok := (*claims)["user_id"].(string)

	if !ok {
		return uuid.Nil, status.Error(codes.InvalidArgument, "token is not valid")
	}

	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "token is not valid")
	}

	return userId, nil
}

func generateTokens(userId uuid.UUID) (string, string, error) {
	accessToken, err := generateAccessToken(userId)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateAccessToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId.String(),
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte("your-256-bit-secret") // TODO replace
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

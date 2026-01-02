package services

import (
	"context"
	"errors"
	"time"

	"github.com/bagussubagja/backend-payment-gateway-go/config"
	"github.com/bagussubagja/backend-payment-gateway-go/internal/models"
	repository "github.com/bagussubagja/backend-payment-gateway-go/internal/repositories"
	"github.com/bagussubagja/backend-payment-gateway-go/internal/utils"
)

type AuthService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	ValidateToken(tokenString string) (uint, error)
	BlacklistToken(ctx context.Context, token string) error
	IsTokenBlacklisted(ctx context.Context, token string) bool
}

type authService struct {
	userRepo     repository.UserRepository
	jwtSecretKey string
	cfg          *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo:     userRepo,
		cfg:          cfg,
		jwtSecretKey: cfg.JWTSecretKey,
	}
}

func (s *authService) Register(req *models.RegisterRequest) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		FullName:    req.FullName,
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		City:        req.City,
		PostalCode:  req.PostalCode,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(user.ID, s.cfg.JWTSecretKey, s.cfg.JWTExpiration)
	if err != nil {
		return nil, err
	}

	loginResponse := &models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:          user.ID,
			FullName:    user.FullName,
			Username:    user.Username,
			Email:       user.Email,
			Address:     user.Address,
			PhoneNumber: user.PhoneNumber,
			City:        user.City,
			PostalCode:  user.PostalCode,
		},
	}
	return loginResponse, nil
}

func (s *authService) ValidateToken(tokenString string) (uint, error) {
	return utils.ValidateToken(tokenString, s.cfg.JWTSecretKey)
}

func (s *authService) BlacklistToken(ctx context.Context, token string) error {
	return config.RedisClient.Set(ctx, token, "true", 24*time.Hour).Err()
}

func (s *authService) IsTokenBlacklisted(ctx context.Context, token string) bool {
	val, _ := config.RedisClient.Get(ctx, token).Result()
	return val == "true"
}

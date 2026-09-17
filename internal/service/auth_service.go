package service

import (
	"errors"

	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/internal/repository"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Refresh(req dto.RefreshRequest) (*dto.LoginResponse, error)
	Me(userUUIDStr string) (*dto.UserResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.Manager) AuthService {
	return &authService{userRepo: userRepo, jwtManager: jwtManager}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, apperror.Internal("gagal memproses password")
	}

	user := &entity.User{Email: req.Email, Password: string(hash)}
	if err := s.userRepo.Create(user); err != nil {
		return nil, apperror.Conflict("email sudah terdaftar")
	}

	return &dto.UserResponse{UUID: user.UUID, Email: user.Email}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("Email atau password salah")
		}
		return nil, apperror.Internal("gagal mengambil data user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.Unauthorized("Email atau password salah")
	}

	return s.issueTokens(user)
}

func (s *authService) Refresh(req dto.RefreshRequest) (*dto.LoginResponse, error) {
	claims, err := s.jwtManager.Parse(req.RefreshToken)
	if err != nil {
		return nil, apperror.Unauthorized("refresh token tidak valid atau kedaluwarsa")
	}

	user, err := s.userRepo.FindByUUID(claims.Sub)
	if err != nil {
		return nil, apperror.Unauthorized("user tidak ditemukan")
	}

	return s.issueTokens(user)
}

func (s *authService) issueTokens(user *entity.User) (*dto.LoginResponse, error) {
	access, err := s.jwtManager.GenerateAccessToken(user.UUID, user.Email)
	if err != nil {
		return nil, apperror.Internal("gagal membuat token")
	}
	refresh, err := s.jwtManager.GenerateRefreshToken(user.UUID, user.Email)
	if err != nil {
		return nil, apperror.Internal("gagal membuat token")
	}

	return &dto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtManager.AccessTTLSeconds(),
		User:         dto.UserResponse{UUID: user.UUID, Email: user.Email},
	}, nil
}

func (s *authService) Me(userUUIDStr string) (*dto.UserResponse, error) {
	id, err := parseUUID(userUUIDStr)
	if err != nil {
		return nil, apperror.Unauthorized("token tidak valid")
	}
	user, err := s.userRepo.FindByUUID(id)
	if err != nil {
		return nil, apperror.Unauthorized("user tidak ditemukan")
	}
	return &dto.UserResponse{UUID: user.UUID, Email: user.Email}, nil
}

package services

import (
	"errors"
	"hackathon/dto"
	"hackathon/models"
	"hackathon/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists   = errors.New("username already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
)

type JwtCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo  repositories.UserRepository
	jwtSecret []byte
	expHours  int
}

func NewAuthService(repo repositories.UserRepository, secret []byte, expHours int) *AuthService {
	return &AuthService{userRepo: repo, jwtSecret: secret, expHours: expHours}
}

func (s *AuthService) Register(username, password string) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{Username: username, Password: string(hashed)}
	if err := s.userRepo.Create(user); err != nil {
		return ErrUserExists
	}
	return nil
}

func (s *AuthService) Login(username, password string) (dto.TokenResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return dto.TokenResponse{}, ErrInvalidCreds
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return dto.TokenResponse{}, ErrInvalidCreds
	}

	exp := time.Now().Add(time.Hour * time.Duration(s.expHours))
	claims := &JwtCustomClaims{
		user.Username,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(exp)},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{Token: token, ExpiredTime: exp.Unix()}, nil
}

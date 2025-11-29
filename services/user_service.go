package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/barzaevhalid/cloud_storage_backend/models"
	"github.com/barzaevhalid/cloud_storage_backend/repositories"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository *repositories.UserRepository
	jwtKey         []byte
	accessExp      time.Duration
	refreshExp     time.Duration
}

func NewUserService(r *repositories.UserRepository, jwtSecret string, accessMin int64, refreshDays int64) *UserService {
	return &UserService{
		UserRepository: r,
		jwtKey:         []byte(jwtSecret),
		accessExp:      time.Minute * time.Duration(accessMin),
		refreshExp:     time.Hour * 24 * time.Duration(refreshDays),
	}
}
func (s *UserService) Register(email, fullname, password string) (*models.User, error) {
	if exist, _ := s.UserRepository.GetByEmail(email); exist != nil && exist.ID != 0 {
		return nil, fmt.Errorf("%s already exists", email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullname,
	}

	if err := s.UserRepository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(email, password string) (string, string, error) {
	user, err := s.UserRepository.Login(email)

	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return "", "", err
	}

	// Createing access token
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(120 * time.Minute).Unix(),
	})
	accessToken, err := at.SignedString(s.jwtKey)

	if err != nil {
		return "", "", err
	}

	//Creating refresh token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	refreshToken, err := rt.SignedString(s.jwtKey)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func (s *UserService) RefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtKey), nil
	})

	if err != nil {
		return "", errors.New("invalid refresh token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalind token claims")
	}

	userId := int64(claims["user_id"].(float64))

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(120 * time.Minute).Unix(),
	})

	return newToken.SignedString(s.jwtKey)
}

func (s *UserService) GetMe(id int64) (*models.User, error) {
	user, err := s.UserRepository.GetMe(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

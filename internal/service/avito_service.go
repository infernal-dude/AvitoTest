package service

import (
	"avito/internal/auth"
	"avito/internal/domain"
	"avito/internal/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

const signKey = "avitoTheBest"

type AvitoService interface {
	Register(user domain.User) error
	Login(password, email string) (string, error)
	LoginDummy(userType string) (string, error)
	CreateHouse(house domain.House) (domain.House, error)
	CreateFlat(flat domain.Flat) (domain.Flat, error)
	UpdateForMod(id int, status string) (domain.Flat, error)
	GetFlatsByHouseId(houseId int, status string) ([]domain.Flat, error)
}

type avitoService struct {
	repository repository.AvitoRepository
}

func NewService(repository repository.AvitoRepository) AvitoService {
	return &avitoService{repository: repository}
}

func (s *avitoService) Register(user domain.User) error {
	var err error
	user.Password, err = auth.GeneratePassword(user.Password)
	if err != nil {
		return err
	}
	err = s.repository.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *avitoService) Login(password, email string) (string, error) {
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("Invalid email or password")
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *avitoService) LoginDummy(userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   uuid.Must(uuid.NewV7()),
		"user_type": userType,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signKey))
}

func (s *avitoService) CreateHouse(house domain.House) (domain.House, error) {
	var houseResp domain.House
	var err error
	houseResp, err = s.repository.CreateHouse(house)
	if err != nil {
		return houseResp, err
	}

	return houseResp, nil
}

func (s *avitoService) CreateFlat(flat domain.Flat) (domain.Flat, error) {
	return s.repository.CreateFlat(flat)
}

func (s *avitoService) UpdateForMod(id int, status string) (domain.Flat, error) {
	return s.repository.UpdateForMod(id, status)
}

func (s *avitoService) GetFlatsByHouseId(houseId int, status string) ([]domain.Flat, error) {
	return s.repository.GetFlatsByHouseId(houseId, status)
}

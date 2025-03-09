package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/user-service/internal/models"
	"github.com/robaa12/keploy-ContractTesting-MicroServices/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser(request models.CreateUserRequest) (models.UserResponse, error)
	GetUserByID(id string) (models.UserResponse, error)
	GetUserByEmail(email string) (models.UserResponse, error)
	ListUsers() ([]models.UserResponse, error)
	UpdateUser(id string, userData models.User) (models.UserResponse, error)
	DeleteUser(id string) error
}

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(request models.CreateUserRequest) (models.UserResponse, error) {
	// Check if user with that email already exists
	_, err := s.repo.GetUserByEmail(request.Email)
	if err == nil {
		return models.UserResponse{}, errors.New("user with that email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, err
	}

	user := models.User{
		ID:        uuid.New().String(),
		Name:      request.Name,
		Email:     request.Email,
		Password:  string(hashedPassword),
		Address:   request.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return models.UserResponse{}, err
	}
	return createdUser.ToUserResponse(), nil
}

func (s *UserService) GetUserByID(id string) (models.UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return models.UserResponse{}, err
	}
	return user.ToUserResponse(), nil
}

func (s *UserService) GetUserByEmail(email string) (models.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return models.UserResponse{}, err
	}
	return user.ToUserResponse(), nil
}

func (s *UserService) ListUsers() ([]models.UserResponse, error) {
	users, err := s.repo.ListUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToUserResponse())
	}
	return userResponses, nil
}

func (s *UserService) UpdateUser(id string, userData models.User) (models.UserResponse, error) {
	// Get the user from the database
	existingUser, err := s.repo.GetUserByID(id)
	if err != nil {
		return models.UserResponse{}, err
	}

	// Update fields
	existingUser.Name = userData.Name
	existingUser.Email = userData.Email
	existingUser.Address = userData.Address
	existingUser.UpdatedAt = time.Now()

	// Save the updated user
	err = s.repo.UpdateUser(existingUser)
	if err != nil {
		return models.UserResponse{}, err
	}

	// Get the updated user
	updatedUser, err := s.repo.GetUserByID(id)
	if err != nil {
		return models.UserResponse{}, err
	}

	return updatedUser.ToUserResponse(), nil
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

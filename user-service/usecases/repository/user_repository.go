package repository

import (
	"user-service/models"
	"user-service/user-proto/users"
)

type UserRepository interface {
	Create(user *models.UserPayload) (int, error)
	FindById(id int64) (*users.UserResponse, error)
	FindUsers() ([]*models.User, error)
	FindByUsername(username string) (*models.User, error)
	DeleteById(id int) error
	Update(user *models.UserPayload) error
}

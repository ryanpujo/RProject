package repository

import (
	"user-service/models"
)

type UserRepository interface {
	Create(user *models.UserPayload) (int64, error)
	FindById(id int64) (*models.User, error)
	FindUsers() ([]*models.User, error)
	FindByUsername(username string) (*models.User, error)
	DeleteById(id int64) error
	Update(user *models.UserPayload) error
}

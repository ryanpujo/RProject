package interactor

import (
	"user-service/models"
	"user-service/usecases/repository"
	"user-service/user-proto/users"

	"golang.org/x/crypto/bcrypt"
)

type userInteractor struct {
	UserRepo repository.UserRepository
}

type UserInteractor interface {
	Create(user *models.UserPayload) (int, error)
	FindById(id int64) (*users.UserResponse, error)
	FindByUsername(username string) (*models.User, error)
	FindUsers() ([]*models.User, error)
	Update(user *models.UserPayload) error
	DeleteById(id int) error
}

func NewUserInteractor(userRepo repository.UserRepository) UserInteractor {
	return &userInteractor{UserRepo: userRepo}
}

func (ui *userInteractor) Create(user *models.UserPayload) (id int, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	user.Password = string(hash)
	id, err = ui.UserRepo.Create(user)
	return
}

func (ui *userInteractor) FindById(id int64) (user *users.UserResponse, err error) {
	user, err = ui.UserRepo.FindById(id)
	return
}

func (ui *userInteractor) FindUsers() (users []*models.User, err error) {
	users, err = ui.UserRepo.FindUsers()
	return
}

func (ui *userInteractor) Update(user *models.UserPayload) (err error) {
	err = ui.UserRepo.Update(user)
	return
}

func (ui *userInteractor) DeleteById(id int) (err error) {
	err = ui.UserRepo.DeleteById(id)
	return
}

func (ui *userInteractor) FindByUsername(username string) (user *models.User, err error) {
	user, err = ui.UserRepo.FindByUsername(username)
	return
}

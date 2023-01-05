package interactor

import (
	"user-service/models"
	"user-service/usecases/repository"

	"golang.org/x/crypto/bcrypt"
)

type userInteractor struct {
	UserRepo repository.UserRepository
}

type UserInteractor interface {
	Create(user *models.UserPayload) (int64, error)
	FindById(id int64) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindUsers() ([]*models.User, error)
	Update(user *models.UserPayload) error
	DeleteById(id int64) error
}

func NewUserInteractor(userRepo repository.UserRepository) *userInteractor {
	return &userInteractor{UserRepo: userRepo}
}

func (ui *userInteractor) Create(user *models.UserPayload) (int64, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)
	id, err := ui.UserRepo.Create(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ui *userInteractor) FindById(id int64) (user *models.User, err error) {
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

func (ui *userInteractor) DeleteById(id int64) (err error) {
	err = ui.UserRepo.DeleteById(id)
	return
}

func (ui *userInteractor) FindByUsername(username string) (user *models.User, err error) {
	user, err = ui.UserRepo.FindByUsername(username)
	return
}

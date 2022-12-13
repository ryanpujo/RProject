package controller

import (
	"context"
	"user-service/models"
	"user-service/usecases/interactor"
	"user-service/user-proto/users"
)

type userServer struct {
	users.UnimplementedUserServiceServer
	Interactor interactor.UserInteractor
}

func NewUserServer(i interactor.UserInteractor) users.UserServiceServer {
	return &userServer{Interactor: i}
}

// create new user if succesfull it will return the id of the newly created user, if not it wil return an error and id equal to 0
func (u *userServer) CreateUser(ctx context.Context, req *users.UserPayload) (*users.Userid, error) {
	input := req.GetUser()
	newUser := models.UserPayload{
		Fname:    input.Fname,
		Lname:    input.Lname,
		Email:    input.Email,
		Username: input.Username,
		Password: req.Password,
	}

	result, err := u.Interactor.Create(&newUser)
	if err != nil {
		return &users.Userid{Id: int64(result)}, err
	}
	res := &users.Userid{Id: int64(result)}
	return res, nil
}

func (u *userServer) FindById(ctx context.Context, id *users.Userid) (user *users.UserResponse, err error) {
	user, err = u.Interactor.FindById(id.GetId())
	return
}

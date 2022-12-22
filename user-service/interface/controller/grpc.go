package controller

import (
	"context"
	"user-service/models"
	"user-service/usecases/interactor"
	"user-service/user-proto/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		return &users.Userid{Id: result}, err
	}
	res := &users.Userid{Id: int64(result)}
	return res, nil
}

func (u *userServer) FindById(ctx context.Context, id *users.Userid) (*users.UserResponse, error) {
	user, err := u.Interactor.FindById(id.GetId())
	if err != nil {
		return nil, err
	}
	result := users.UserResponse{
		User: &users.User{
			Id:        int64(user.Id),
			Fname:     user.Fname,
			Lname:     user.Lname,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}
	return &result, nil
}

func (u *userServer) FindUsers(empty *emptypb.Empty, stream users.UserService_FindUsersServer) error {
	data, err := u.Interactor.FindUsers()
	if err != nil {
		return status.Errorf(codes.DataLoss, "failed to retrieve error")
	}
	for _, v := range data {
		user := users.UserResponse{
			User: &users.User{
				Id:        int64(v.Id),
				Fname:     v.Fname,
				Lname:     v.Lname,
				Email:     v.Email,
				Username:  v.Username,
				CreatedAt: timestamppb.New(v.CreatedAt),
				UpdatedAt: timestamppb.New(v.UpdatedAt),
			},
		}
		err = stream.Send(&user)
		if err != nil {
			return status.Errorf(codes.DataLoss, "failed to retrieve error")
		}
	}
	return nil
}

package controller_test

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"testing"
	"time"
	"user-service/interface/controller"
	"user-service/models"
	"user-service/user-proto/users"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mockInteractor struct {
	mock.Mock
}

func (m *mockInteractor) Create(user *models.UserPayload) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockInteractor) FindById(id int64) (*models.User, error) {
	args := m.Called(id)
	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}
	return user, args.Error(1)
}

func (m *mockInteractor) FindUsers() ([]*models.User, error) {
	args := m.Called()
	var users []*models.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*models.User)
	}
	return users, args.Error(1)
}

func (m *mockInteractor) Update(user *models.UserPayload) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockInteractor) DeleteById(id int64) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockInteractor) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}
	return user, args.Error(1)
}

var mocking *mockInteractor
var lis *bufconn.Listener
var client users.UserServiceClient

func buffdialer(ctx context.Context, s string) (net.Conn, error) {
	return lis.Dial()
}
func TestMain(m *testing.M) {
	lis = bufconn.Listen(1024 * 1024)
	defer lis.Close()
	s := grpc.NewServer()
	defer s.Stop()
	mocking = new(mockInteractor)
	users.RegisterUserServiceServer(s, controller.NewUserServer(mocking))
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(buffdialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial microservice: %s", err)
	}
	defer conn.Close()
	client = users.NewUserServiceClient(conn)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	os.Exit(m.Run())
}

func Test_userServer_CreateUser(t *testing.T) {
	user := &users.User{
		Fname:    "ryan",
		Lname:    "pujo",
		Username: "ryanpujo",
		Email:    "ryanpujo@yahoo.com",
	}
	jsonReq := &users.UserPayload{
		Password: "supersecret1",
		User:     user,
	}
	testTable := map[string]struct {
		jsonReq *users.UserPayload
		arrange func(t *testing.T)
		assert  func(t *testing.T, e error, id int64)
	}{
		"success api call": {
			jsonReq: jsonReq,
			arrange: func(t *testing.T) {
				mocking.On("Create", mock.Anything).Return(int64(1), nil).Once()
			},
			assert: func(t *testing.T, e error, id int64) {
				require.Equal(t, int64(1), id)
				require.NoError(t, e)
			},
		},
		"fail api call": {
			jsonReq: jsonReq,
			arrange: func(t *testing.T) {
				mocking.On("Create", mock.Anything).Return(int64(0), errors.New("got error")).Once()
			},
			assert: func(t *testing.T, e error, id int64) {
				require.Zero(t, id)
				require.Error(t, e)
			},
		},
	}

	for k, v := range testTable {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		t.Run(k, func(t *testing.T) {
			v.arrange(t)

			res, err := client.CreateUser(ctx, v.jsonReq)

			v.assert(t, err, res.GetId())
		})
	}
}

func Test_userServer_FindById(t *testing.T) {
	user := &models.User{
		Id:        1,
		Fname:     "ryan",
		Lname:     "pujo",
		Username:  "ryanpujo",
		Email:     "ryanpujo@yahoo.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id := users.Userid{Id: int64(1)}
	testTable := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, e error, actual *users.UserResponse)
	}{
		"succes api call": {
			arrange: func(t *testing.T) {
				mocking.On("FindById", int64(1)).Return(user, nil).Once()
			},
			assert: func(t *testing.T, e error, actual *users.UserResponse) {
				require.NoError(t, e)
				require.Equal(t, user.Fname, actual.GetUser().GetFname())
			},
		},
		"failed api call": {
			arrange: func(t *testing.T) {
				mocking.On("FindById", int64(1)).Return(nil, errors.New("user not found")).Once()
			},
			assert: func(t *testing.T, e error, actual *users.UserResponse) {
				require.Error(t, e)
				require.Empty(t, actual)
			},
		},
	}

	for k, v := range testTable {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		t.Run(k, func(t *testing.T) {
			v.arrange(t)

			res, err := client.FindById(ctx, &id)

			v.assert(t, err, res)
		})
	}
}

func Test_userServer_FindUsers(t *testing.T) {
	user := &models.User{
		Id:        1,
		Fname:     "ryan",
		Lname:     "pujo",
		Username:  "ryanpujo",
		Email:     "ryanpujo@yahoo.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user1 := &models.User{
		Id:        1,
		Fname:     "ryan",
		Lname:     "pujo",
		Username:  "ryanpujo",
		Email:     "ryanpujo@yahoo.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	allUsers := []*models.User{user, user1}
	testTable := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, err error, actual users.UserService_FindUsersClient)
	}{
		"successful call": {
			arrange: func(t *testing.T) {
				mocking.On("FindUsers").Return(allUsers, nil).Once()
			},
			assert: func(t *testing.T, err error, actual users.UserService_FindUsersClient) {
				require.NoError(t, err)
				require.NotNil(t, actual)
				res, _ := actual.Recv()
				require.Equal(t, res.GetUser().Fname, user.Fname)
			},
		},
		"failed call": {
			arrange: func(t *testing.T) {
				mocking.On("FindUsers").Return(nil, status.Error(codes.DataLoss, "no user found")).Once()
			},
			assert: func(t *testing.T, err error, actual users.UserService_FindUsersClient) {
				res, _ := actual.Recv()
				require.Empty(t, res.GetUser().GetFname())
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.arrange(t)

			result, err := client.FindUsers(ctx, &emptypb.Empty{})

			v.assert(t, err, result)
		})
	}
}

func Test_userServer_Update(t *testing.T) {
	user := users.User{
		Id:        1,
		Fname:     "ryan",
		Lname:     "pujo",
		Username:  "ryanpujo",
		Email:     "ryanpujo@yahoo.com",
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}
	payload := users.UserPayload{
		User:     &user,
		Password: "ksjdjnejn",
	}
	testTable := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, err error)
	}{
		"success call": {
			arrange: func(t *testing.T) {
				mocking.On("Update", mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed call": {
			arrange: func(t *testing.T) {
				mocking.On("Update", mock.Anything).Return(errors.New("got an error")).Once()
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "rpc error: code = InvalidArgument desc = got an error", err.Error())
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.arrange(t)

			_, err := client.Update(ctx, &payload)

			v.assert(t, err)
		})
	}
}

func Test_userServer_DeleteById(t *testing.T) {
	testTable := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, err error)
	}{
		"success call": {
			arrange: func(t *testing.T) {
				mocking.On("DeleteById", mock.Anything).Return(nil).Once()
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed call": {
			arrange: func(t *testing.T) {
				mocking.On("DeleteById", mock.Anything).Return(errors.New("got an error")).Once()
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, "rpc error: code = InvalidArgument desc = got an error", err.Error())
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.arrange(t)

			_, err := client.DeleteById(ctx, &users.Userid{Id: 1})

			v.assert(t, err)
		})
	}
}

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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type mockInteractor struct {
	mock.Mock
}

func (m *mockInteractor) Create(user *models.UserPayload) (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *mockInteractor) FindById(id int64) (*users.UserResponse, error) {
	args := m.Called(id)
	var user *users.UserResponse
	if args.Get(0) != nil {
		user = args.Get(0).(*users.UserResponse)
	}
	return user, args.Error(1)
}

func (m *mockInteractor) FindUsers() ([]*models.User, error) {
	args := m.Called()
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *mockInteractor) Update(user *models.UserPayload) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockInteractor) DeleteById(id int) error {
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
	s := grpc.NewServer()
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
				mocking.On("Create", mock.Anything).Return(1, nil).Once()
			},
			assert: func(t *testing.T, e error, id int64) {
				require.Equal(t, int64(1), id)
				require.NoError(t, e)
			},
		},
		"fail api call": {
			jsonReq: jsonReq,
			arrange: func(t *testing.T) {
				mocking.On("Create", mock.Anything).Return(0, errors.New("got error")).Once()
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
	user := &users.User{
		Fname:    "ryan",
		Lname:    "pujo",
		Username: "ryanpujo",
		Email:    "ryanpujo@yahoo.com",
	}
	id := users.Userid{Id: 1}
	res := &users.UserResponse{User: user}
	testTable := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, e error, actual *users.UserResponse)
	}{
		"succes api call": {
			arrange: func(t *testing.T) {
				mocking.On("FindById", int64(1)).Return(res, nil).Once()
			},
			assert: func(t *testing.T, e error, actual *users.UserResponse) {
				require.NoError(t, e)
				require.Equal(t, res.GetUser().GetFname(), actual.GetUser().GetFname())
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

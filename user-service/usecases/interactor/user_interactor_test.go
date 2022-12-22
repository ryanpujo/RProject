package interactor_test

import (
	"errors"
	"os"
	"testing"
	"user-service/models"
	"user-service/usecases/interactor"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	mock.Mock
}

func (ui *mockUserRepo) Create(user *models.UserPayload) (int64, error) {
	args := ui.Called(user)
	return args.Get(0).(int64), args.Error(1)
}

func (ui *mockUserRepo) FindById(id int64) (user *models.User, err error) {
	args := ui.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (ui *mockUserRepo) FindUsers() (users []*models.User, err error) {
	users = []*models.User{
		{
			Fname:    "ryan",
			Lname:    "pujo",
			Username: "ryanpujo",
			Email:    "ryanpujo",
			Password: "secret",
		},
	}

	return
}

func (ui *mockUserRepo) Update(user *models.UserPayload) (err error) {
	return
}

func (ui *mockUserRepo) DeleteById(id int64) (err error) {
	return
}

func (ur *mockUserRepo) FindByUsername(username string) (user *models.User, err error) {
	args := ur.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

var userInteractorTest interactor.UserInteractor
var mockRepo *mockUserRepo

func TestMain(m *testing.M) {
	mockRepo = new(mockUserRepo)
	userInteractorTest = interactor.NewUserInteractor(mockRepo)
	os.Exit(m.Run())
}

func Test_userInteractor_Create(t *testing.T) {
	user := models.UserPayload{
		Fname:    "ryan",
		Lname:    "pujo",
		Username: "ryanpujo",
		Email:    "ryanpujo",
		Password: "secret",
	}
	testTable := map[string]struct {
		payload models.UserPayload
		arrange func(t *testing.T)
		assert  func(t *testing.T, err error, id int64)
	}{
		"success create": {
			payload: user,
			arrange: func(t *testing.T) {
				mockRepo.On("Create", mock.Anything).Return(int64(1), nil).Once()
			},
			assert: func(t *testing.T, err error, id int64) {
				require.NoError(t, err)
				require.Equal(t, int64(1), id)
			},
		},
		"fail create": {
			payload: user,
			arrange: func(t *testing.T) {
				mockRepo.On("Create", mock.Anything).Return(int64(0), errors.New("got error")).Once()
			},
			assert: func(t *testing.T, err error, id int64) {
				require.Error(t, err)
				require.Equal(t, int64(0), id)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.arrange(t)

			result, err := userInteractorTest.Create(&v.payload)

			v.assert(t, err, int64(result))
		})
	}
}

func Test_userInteractor_FindById(t *testing.T) {
	userTest := models.User{
		Fname:    "ryan",
		Lname:    "pujo",
		Username: "ryanpujo",
		Email:    "ryanpujo",
	}

	mockRepo.On("FindById", int64(1)).Return(&userTest, nil).Once()
	user, err := userInteractorTest.FindById(int64(1))
	if err != nil {
		t.Errorf("failed to retrieve user by id 1: %s", err)
	}

	if user.Fname != userTest.Fname {
		t.Errorf("wrong name; expect %s but got %s", userTest.Fname, user.Fname)
	}
}

func Test_userInteractor_FindByUsername(t *testing.T) {
	userTest := models.User{
		Fname:    "ryan",
		Lname:    "pujo",
		Username: "ryanpujo",
		Email:    "ryanpujo",
		Password: "secret",
	}
	mockRepo.On("FindByUsername", "ryanpujo").Return(&userTest, nil).Once()
	user, err := userInteractorTest.FindByUsername("ryanpujo")
	if err != nil {
		t.Errorf("failed to retrieve user by username ryanpujo: %s", err)
	}

	if user.Username != userTest.Username {
		t.Errorf("wrong name; expect %s but got %s", userTest.Username, user.Username)
	}
}

func Test_userInteractor_FindUsers(t *testing.T) {
	users, err := userInteractorTest.FindUsers()
	if err != nil {
		t.Errorf("failed to retrieve users: %s", err)
	}

	if len(users) != 1 {
		t.Errorf("wrong number of users; expect 1 but got %d", len(users))
	}
}

func Test_userInteractor_Update(t *testing.T) {
	user := models.UserPayload{
		Fname:    "ryan",
		Lname:    "pujo",
		Username: "ryanpujo",
		Email:    "ryanpujo@yahoo.com",
		Password: "supersecret1",
	}
	user.Lname = "conor"
	err := userInteractorTest.Update(&user)
	if err != nil {
		t.Errorf("failed to update user: %s", err)
	}
}

func Test_userInteractor_DeleteById(t *testing.T) {
	err := userInteractorTest.DeleteById(1)
	if err != nil {
		t.Errorf("failed to delete a user by id 1: %s", err)
	}
}

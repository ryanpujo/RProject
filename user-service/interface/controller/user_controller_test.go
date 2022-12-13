package controller

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
// 	"user-service/interface/controller"
// 	"user-service/models"
// 	"user-service/router"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// type mockInteractor struct {
// 	mock.Mock
// }

// func (m *mockInteractor) Create(user *models.UserPayload) (int, error) {
// 	args := m.Called()
// 	return args.Get(0).(int), args.Error(1)
// }

// func (m *mockInteractor) FindById(id int) (*models.User, error) {
// 	args := m.Called(id)
// 	var user *models.User
// 	if args.Get(0) != nil {
// 		user = args.Get(0).(*models.User)
// 	}
// 	return user, args.Error(1)
// }

// func (m *mockInteractor) FindUsers() ([]*models.User, error) {
// 	args := m.Called()
// 	return args.Get(0).([]*models.User), args.Error(1)
// }

// func (m *mockInteractor) Update(user *models.UserPayload) error {
// 	args := m.Called()
// 	return args.Error(0)
// }

// func (m *mockInteractor) DeleteById(id int) error {
// 	args := m.Called()
// 	return args.Error(0)
// }

// func (m *mockInteractor) FindByUsername(username string) (*models.User, error) {
// 	args := m.Called(username)
// 	var user *models.User
// 	if args.Get(0) != nil {
// 		user = args.Get(0).(*models.User)
// 	}
// 	return user, args.Error(1)
// }

// var userController controller.UserController
// var mocking *mockInteractor
// var route *gin.Engine

// func TestMain(m *testing.M) {
// 	gin.SetMode(gin.TestMode)
// 	mocking = new(mockInteractor)
// 	userController = controller.NewUserController(mocking)
// 	appController := controller.AppController{
// 		User: userController,
// 	}
// 	route = router.NewRouter(appController)
// 	os.Exit(m.Run())
// }

// type testCreateResponse struct {
// 	Error   bool   `json:"error"`
// 	Message string `json:"message"`
// 	Data    int    `json:"data,omitempty"`
// }

// func Test_userController_Create(t *testing.T) {
// 	jsonReq := []byte(`{
// 			"fname":    "ryan",
// 			"lname":    "pujo",
// 			"username": "ryanpujo",
// 			"email":    "ryanpujo@yahoo.com",
// 			"password": "supersecret1"
// 		}`)
// 	badJson := []byte(`{
// 		"fname":    "ryan",
// 		"lname":    "pujo",
// 		"username": "ryanpujo",
// 		"email":    "ryanpujo@yahoo.com",
// 		"password": "supersecret"
// 	}`)
// 	testTable := map[string]struct {
// 		jsonReq []byte
// 		arrange func(t *testing.T)
// 		assert  func(t *testing.T, statusCode int, err bool, message ...string)
// 	}{
// 		"success api call": {
// 			jsonReq: jsonReq,
// 			arrange: func(t *testing.T) {
// 				mocking.On("Create", mock.Anything).Return(1, nil).Once()
// 			},
// 			assert: func(t *testing.T, statusCode int, err bool, message ...string) {
// 				if len(message) > 0 {
// 					require.Equal(t, "USER CREATED", message[0])
// 				}
// 				require.Equal(t, http.StatusCreated, statusCode)
// 				require.False(t, err)
// 			},
// 		},
// 		"fail api call": {
// 			jsonReq: jsonReq,
// 			arrange: func(t *testing.T) {
// 				mocking.On("Create", mock.Anything).Return(0, errors.New("got error")).Once()
// 			},
// 			assert: func(t *testing.T, statusCode int, err bool, message ...string) {
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				require.True(t, err)
// 				if len(message) > 0 {
// 					require.Equal(t, "got error", message[0])
// 				}
// 			},
// 		},
// 		"bad json": {
// 			jsonReq: badJson,
// 			arrange: func(t *testing.T) {

// 			},
// 			assert: func(t *testing.T, statusCode int, err bool, message ...string) {
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				require.True(t, err)
// 				if len(message) > 0 {
// 					require.Equal(t, "Key: 'UserPayload.Password' Error:Field validation for 'Password' failed on the 'min' tag", message[0])
// 				}
// 			},
// 		},
// 	}

// 	for k, v := range testTable {
// 		t.Run(k, func(t *testing.T) {
// 			v.arrange(t)

// 			rr := httptest.NewRecorder()
// 			req, err := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(v.jsonReq))
// 			require.Nil(t, err)
// 			req.Header.Set("Content-Type", "application/json")
// 			route.ServeHTTP(rr, req)
// 			var response testCreateResponse
// 			_ = json.NewDecoder(rr.Body).Decode(&response)

// 			v.assert(t, rr.Code, response.Error, response.Message)
// 		})
// 	}
// }

// func Test_userController_FindById(t *testing.T) {
// 	user := models.User{
// 		Id:       1,
// 		Fname:    "ryan",
// 		Lname:    "pujo",
// 		Username: "ryanpujo",
// 		Email:    "ryanpujo",
// 	}
// 	testTable := map[string]struct {
// 		url     string
// 		arrange func(t *testing.T)
// 		assert  func(t *testing.T, actual models.User, err bool, statusCode int, messages ...string)
// 	}{
// 		"success call": {
// 			url: "/user/1",
// 			arrange: func(t *testing.T) {
// 				mocking.On("FindById", 1).Return(&user, nil).Once()
// 			},
// 			assert: func(t *testing.T, actual models.User, err bool, statusCode int, messages ...string) {
// 				require.Equal(t, user, actual)
// 				require.False(t, err)
// 				require.Equal(t, http.StatusOK, statusCode)
// 				if len(messages) > 0 {
// 					require.Equal(t, "FOUND USER", messages[0])
// 				}
// 			},
// 		},
// 		"fail in interactor": {
// 			url: "/user/1",
// 			arrange: func(t *testing.T) {
// 				mocking.On("FindById", 1).Return(nil, errors.New("user not found")).Once()
// 			},
// 			assert: func(t *testing.T, actual models.User, err bool, statusCode int, messages ...string) {
// 				require.Empty(t, actual.Fname)
// 				require.True(t, err)
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				if len(messages) > 0 {
// 					require.Equal(t, "user not found", messages[0])
// 				}
// 			},
// 		},
// 		"bad params": {
// 			url: "/user/-1",
// 			arrange: func(t *testing.T) {
// 			},
// 			assert: func(t *testing.T, actual models.User, err bool, statusCode int, messages ...string) {
// 				require.Empty(t, actual.Fname)
// 				require.True(t, err)
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				if len(messages) > 0 {
// 					require.Equal(t, "Key: 'Uri.ID' Error:Field validation for 'ID' failed on the 'min' tag", messages[0])
// 				}
// 			},
// 		},
// 	}

// 	for k, v := range testTable {
// 		t.Run(k, func(t *testing.T) {
// 			v.arrange(t)

// 			req := httptest.NewRequest("GET", v.url, nil)
// 			rr := httptest.NewRecorder()
// 			route.ServeHTTP(rr, req)
// 			var res struct {
// 				Error   bool        `json:"error"`
// 				Message string      `json:"message"`
// 				Data    models.User `json:"data,omitempty"`
// 			}
// 			json.NewDecoder(rr.Body).Decode(&res)

// 			v.assert(t, res.Data, res.Error, rr.Code, res.Message)
// 		})
// 	}
// }

// func Test_userController_FindUsers(t *testing.T) {
// 	user1 := models.User{
// 		Id:       1,
// 		Fname:    "ryan",
// 		Lname:    "pujo",
// 		Username: "ryanpujo",
// 		Email:    "ryanpujo",
// 	}
// 	user2 := models.User{
// 		Id:       1,
// 		Fname:    "ryan",
// 		Lname:    "pujo",
// 		Username: "ryanpujo",
// 		Email:    "ryanpujo",
// 	}
// 	user3 := models.User{
// 		Id:       1,
// 		Fname:    "ryan",
// 		Lname:    "pujo",
// 		Username: "ryanpujo",
// 		Email:    "ryanpujo",
// 	}
// 	empty := []*models.User(nil)
// 	users := []*models.User{&user1, &user2, &user3}
// 	testTable := map[string]struct {
// 		arrange func(t *testing.T)
// 		assert  func(t *testing.T, actual []*models.User, n, statusCode int, err bool, messages ...string)
// 	}{
// 		"succes call": {
// 			arrange: func(t *testing.T) {
// 				mocking.On("FindUsers").Return(users, nil).Once()
// 			},
// 			assert: func(t *testing.T, actual []*models.User, n, statusCode int, err bool, messages ...string) {
// 				require.Equal(t, users, actual)
// 				require.Equal(t, 3, n)
// 				require.Equal(t, http.StatusOK, statusCode)
// 				require.False(t, err)
// 				if len(messages) > 0 {
// 					require.Equal(t, "FOUND USERS", messages[0])
// 				}
// 			},
// 		},
// 		"fail call": {
// 			arrange: func(t *testing.T) {
// 				mocking.On("FindUsers").Return(empty, errors.New("no users found")).Once()
// 			},
// 			assert: func(t *testing.T, actual []*models.User, n, statusCode int, err bool, messages ...string) {
// 				require.Equal(t, empty, actual)
// 				require.Equal(t, 0, n)
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				require.True(t, err)
// 				if len(messages) > 0 {
// 					require.Equal(t, "no users found", messages[0])
// 				}
// 			},
// 		},
// 	}

// 	for k, v := range testTable {
// 		t.Run(k, func(t *testing.T) {
// 			v.arrange(t)

// 			req := httptest.NewRequest("GET", "/user", nil)
// 			rr := httptest.NewRecorder()
// 			route.ServeHTTP(rr, req)
// 			var res struct {
// 				Error   bool           `json:"error"`
// 				Message string         `json:"message"`
// 				Data    []*models.User `json:"data,omitempty"`
// 			}
// 			json.NewDecoder(rr.Body).Decode(&res)

// 			v.assert(t, res.Data, len(res.Data), rr.Code, res.Error, res.Message)
// 		})
// 	}
// }

// func Test_userController_Update(t *testing.T) {
// 	jsonReq, _ := json.Marshal(models.UserPayload{
// 		Fname:    "ryan",
// 		Lname:    "pujo",
// 		Username: "ryanpujo",
// 		Email:    "ryanpujo@yahoo.com",
// 		Password: "supersecret1",
// 	})
// 	badJson, _ := json.Marshal(models.UserPayload{
// 		Fname:    "ryan",
// 		Lname:    "pujo",
// 		Username: "ryanpujo",
// 		Email:    "ryanpujo@yahoo.com",
// 		Password: "supersecret",
// 	})
// 	testTable := map[string]struct {
// 		jsonReq []byte
// 		arrange func(t *testing.T)
// 		assert  func(t *testing.T, statusCode int, err bool, messages ...string)
// 	}{
// 		"success api call": {
// 			jsonReq: jsonReq,
// 			arrange: func(t *testing.T) {
// 				mocking.On("Update", mock.Anything).Return(nil).Once()
// 			},
// 			assert: func(t *testing.T, statusCode int, err bool, messages ...string) {
// 				require.Equal(t, http.StatusOK, statusCode)
// 				require.False(t, err)
// 				if len(messages) > 0 {
// 					require.Equal(t, "ACCOUNT UPDATED", messages[0])
// 				}
// 			},
// 		},
// 		"fail api call": {
// 			jsonReq: jsonReq,
// 			arrange: func(t *testing.T) {
// 				mocking.On("Update", mock.Anything).Return(errors.New("fail to update account")).Once()
// 			},
// 			assert: func(t *testing.T, statusCode int, err bool, messages ...string) {
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				require.True(t, err)
// 				if len(messages) > 0 {
// 					require.Equal(t, "fail to update account", messages[0])
// 				}
// 			},
// 		},
// 		"bad json": {
// 			jsonReq: badJson,
// 			arrange: func(t *testing.T) {},
// 			assert: func(t *testing.T, statusCode int, err bool, messages ...string) {
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				require.True(t, err)
// 				if len(messages) > 0 {
// 					require.Equal(t, "Key: 'UserPayload.Password' Error:Field validation for 'Password' failed on the 'min' tag", messages[0])
// 				}
// 			},
// 		},
// 	}

// 	for k, v := range testTable {
// 		t.Run(k, func(t *testing.T) {
// 			v.arrange(t)

// 			req := httptest.NewRequest("PATCH", "/user/update", bytes.NewReader(v.jsonReq))
// 			rr := httptest.NewRecorder()
// 			route.ServeHTTP(rr, req)
// 			var res struct {
// 				Error   bool           `json:"error"`
// 				Message string         `json:"message"`
// 				Data    []*models.User `json:"data,omitempty"`
// 			}
// 			json.NewDecoder(rr.Body).Decode(&res)

// 			v.assert(t, rr.Code, res.Error, res.Message)
// 		})
// 	}
// }

// func Test_userController_DeleteById(t *testing.T) {
// 	testTable := map[string]struct {
// 		urls    string
// 		arrange func(t *testing.T)
// 		assert  func(t *testing.T, statusCode int, err bool, messages ...string)
// 	}{
// 		"success api call": {
// 			urls: "/user/delete/1",
// 			arrange: func(t *testing.T) {
// 				mocking.On("DeleteById", mock.Anything).Return(nil).Once()
// 			},
// 			assert: func(t *testing.T, statusCode int, err bool, messages ...string) {
// 				require.Equal(t, http.StatusOK, statusCode)
// 				require.False(t, err)
// 				if len(messages) > 0 {
// 					require.Equal(t, "ACCOUNT DELETED", messages[0])
// 				}
// 			},
// 		},
// 		"fail in interactor": {
// 			urls: "/user/delete/1",
// 			arrange: func(t *testing.T) {
// 				mocking.On("DeleteById", mock.Anything).Return(errors.New("fail to delete account")).Once()
// 			},
// 			assert: func(t *testing.T, statusCode int, err bool, messages ...string) {
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				require.True(t, err)
// 				if len(messages) > 0 {
// 					require.Equal(t, "fail to delete account", messages[0])
// 				}
// 			},
// 		},
// 		"bad parameter req": {
// 			urls: "/user/delete/-1",
// 			arrange: func(t *testing.T) {

// 			},
// 			assert: func(t *testing.T, statusCode int, err bool, messages ...string) {
// 				require.Equal(t, http.StatusBadRequest, statusCode)
// 				require.True(t, err)
// 				if len(messages) > 0 {
// 					require.Equal(t, "Key: 'Uri.ID' Error:Field validation for 'ID' failed on the 'min' tag", messages[0])
// 				}
// 			},
// 		},
// 	}

// 	for k, v := range testTable {
// 		t.Run(k, func(t *testing.T) {
// 			v.arrange(t)

// 			req := httptest.NewRequest("DELETE", v.urls, nil)
// 			rr := httptest.NewRecorder()
// 			route.ServeHTTP(rr, req)
// 			var res struct {
// 				Error   bool           `json:"error"`
// 				Message string         `json:"message"`
// 				Data    []*models.User `json:"data,omitempty"`
// 			}

// 			json.NewDecoder(rr.Body).Decode(&res)

// 			v.assert(t, rr.Code, res.Error, res.Message)
// 		})
// 	}
// }

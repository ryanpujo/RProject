package controller

// import (
// 	"net/http"
// 	"user-service/models"
// 	"user-service/usecases/interactor"

// 	"github.com/gin-gonic/gin"
// )

// type UserController interface {
// 	Create(c *gin.Context)
// 	FindById(c *gin.Context)
// 	FindUsers(c *gin.Context)
// 	Update(c *gin.Context)
// 	DeleteById(c *gin.Context)
// }

// type userController struct {
// 	userInteractor interactor.UserInteractor
// }

// type JsonResponse struct {
// 	Error   bool   `json:"error"`
// 	Message string `json:"message"`
// 	Data    any    `json:"data,omitempty"`
// }

// type Uri struct {
// 	ID int `uri:"id" binding:"required,min=1"`
// }

// func NewUserController(userInteractor interactor.UserInteractor) UserController {
// 	return &userController{userInteractor: userInteractor}
// }

// func (uc userController) Create(c *gin.Context) {
// 	var user models.UserPayload
// 	var json JsonResponse
// 	err := c.ShouldBindJSON(&user)
// 	if err != nil {
// 		json.Error = true
// 		json.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, json)
// 		return
// 	}
// 	id, errs := uc.userInteractor.Create(&user)
// 	if errs != nil {
// 		json.Error = true
// 		json.Message = errs.Error()
// 		c.JSON(http.StatusBadRequest, json)
// 		return
// 	}

// 	json.Error = false
// 	json.Message = "USER CREATED"
// 	json.Data = id
// 	c.JSON(http.StatusCreated, json)
// }

// func (uc userController) FindById(c *gin.Context) {
// 	var uri Uri
// 	var res JsonResponse
// 	err := c.ShouldBindUri(&uri)
// 	if err != nil {
// 		res.Error = true
// 		res.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// 	user, err := uc.userInteractor.FindById(uri.ID)
// 	if err != nil {
// 		res.Error = true
// 		res.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// 	json := JsonResponse{
// 		Error:   false,
// 		Message: "FOUND USER",
// 		Data:    user,
// 	}
// 	c.JSON(http.StatusOK, json)
// }

// func (uc userController) FindUsers(c *gin.Context) {
// 	users, err := uc.userInteractor.FindUsers()
// 	var res JsonResponse
// 	if err != nil {
// 		res.Error = true
// 		res.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// 	json := JsonResponse{
// 		Error:   false,
// 		Message: "FOUND USERS",
// 		Data:    users,
// 	}
// 	c.JSON(http.StatusOK, json)
// }

// func (uc userController) Update(c *gin.Context) {
// 	var user models.UserPayload
// 	var res JsonResponse
// 	err := c.ShouldBindJSON(&user)
// 	if err != nil {
// 		res.Error = true
// 		res.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// 	err = uc.userInteractor.Update(&user)
// 	if err != nil {
// 		res.Error = true
// 		res.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	json := JsonResponse{
// 		Error:   false,
// 		Message: "ACCOUNT UPDATED",
// 	}

// 	c.JSON(http.StatusOK, json)
// }

// func (uc userController) DeleteById(c *gin.Context) {
// 	var uri Uri
// 	var res JsonResponse
// 	err := c.ShouldBindUri(&uri)
// 	if err != nil {
// 		res.Error = true
// 		res.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// 	err = uc.userInteractor.DeleteById(uri.ID)
// 	if err != nil {
// 		res.Error = true
// 		res.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// 	json := JsonResponse{
// 		Error:   false,
// 		Message: "ACCOUNT DELETED",
// 	}

// 	c.JSON(http.StatusOK, json)
// }

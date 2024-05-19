package controllers

import (
	"banking_application/api/domain/dtos"
	"banking_application/api/services"
	"banking_application/api/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(service services.UserService) UserController {
	return UserController{
		userService: &service,
	}
}

func (s *UserController) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := dtos.UserDto{}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
			return
		}

		_, serviceErr := s.userService.CreateUser(user)

		if serviceErr != nil {
			errResponse := util.HandleErrors(serviceErr)
			c.JSON(errResponse.StatusCode, gin.H{"error": errResponse.Message, "success": errResponse.Success})
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "success": true})
	}
}

func (s *UserController) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := s.userService.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching users"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func (s *UserController) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching users"})
			return
		}

		result, err := s.userService.GetUser(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching users"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

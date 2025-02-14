package controller

import (
	"Petstore/models/user/entity"
	"Petstore/models/user/service"
	"Petstore/responder"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ptflp/godecoder"
	"net/http"
)

type UserController struct {
	service service.UserServicer
	responder.Responder
	godecoder.Decoder
}

func NewUserController(service service.UserServicer, responder responder.Responder, decoder godecoder.Decoder) *UserController {
	return &UserController{service: service, Responder: responder, Decoder: decoder}
}

// Login
// @Summary User login
// @Description Authenticate a user by username and password
// @Tags user
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {string} string "Access token"
// @Router /user/login [post]
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	auth := c.service.Login(r.Context(), service.LoginIn{Username: login, Password: password})
	if auth.Err != nil {
		c.ErrorUnauthorized(w, auth.Err)
		return
	}
	accessToken := auth.AccessToken

	c.OutputJSON(w, accessToken)

}

// Logout
// @Summary User logout
// @Description Terminate the user session
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} LogoutResponse "Logout success"
// @Router /user/logout [post]
func (c *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	err := c.service.Logout(r.Context())
	if err != nil {
		c.OutputJSON(w, Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    "Logout error: " + err.Error(),
		})
		return
	}
	c.OutputJSON(w, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    "Logout success",
	})
}

// GetUser
// @Summary Get user information
// @Description Retrieve user data by username
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} GetUserResponse "User data"
// @Router /user/{username} [get]
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := c.service.GetUser(r.Context(), username)
	if err != nil {
		c.OutputJSON(w, Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    "Error: " + err.Error(),
		})
		return
	}
	c.OutputJSON(w, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    user,
	})
}

// UpdateUser
// @Summary Update user information
// @Description Update user data by username
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param body body entity.User true "User update data"
// @Success 200 {object} GetUserResponse "Updated user data"
// @Router /user/{username} [put]
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	var updateData entity.User
	if err := c.Decode(r.Body, &updateData); err != nil {
		c.ErrorBadRequest(w, err)
		return
	}

	updatedUser, err := c.service.UpdateUser(r.Context(), username, &updateData)
	if err != nil {
		c.OutputJSON(w, Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    "Error: " + err.Error(),
		})
		return
	}

	c.OutputJSON(w, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    updatedUser,
	})
}

// DeleteUser
// @Summary Delete a user
// @Description Delete a user by username
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} LogoutResponse  "User deleted"
// @Router /user/{username} [delete]
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	err := c.service.DeleteUser(r.Context(), username)
	if err != nil {
		c.OutputJSON(w, Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    "Error: " + err.Error(),
		})
		return
	}
	c.OutputJSON(w, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    "User deleted",
	})
}

// CreateUser
// @Summary Create a new user
// @Description Create a new user with the provided data
// @Tags user
// @Accept json
// @Produce json
// @Param body body entity.User true "User data"
// @Success 200 {object} LogoutResponse "User created"
// @Router /user [post]
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := c.Decode(r.Body, &user); err != nil {
		c.ErrorBadRequest(w, err)
		return
	}
	id, err := c.service.CreateUser(r.Context(), &user)
	if err != nil {
		c.OutputJSON(w, Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    "Error: " + err.Error(),
		})
		return
	}
	c.OutputJSON(w, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    fmt.Sprintf("User created: %d", id),
	})

}

// CreateUsersWithArray
// @Summary Create users with an array
// @Description Create multiple users with an array of user data
// @Tags user
// @Accept json
// @Produce json
// @Param body body []entity.User true "Array of user data"
// @Success 200 {object} GetUsersResponse "Users created"
// @Router /user/createWithArray [post]
func (c *UserController) CreateUsersWithArray(w http.ResponseWriter, r *http.Request) {
	var users []*entity.User
	if err := c.Decode(r.Body, &users); err != nil {
		c.ErrorBadRequest(w, err)
		return
	}
	err := c.service.CreateWithArray(r.Context(), users)
	if err != nil {
		c.OutputJSON(w, Response{
			Success: false,
			Code:    http.StatusInternalServerError,
			Data:    "Error: " + err.Error(),
		})
		return
	}
	c.OutputJSON(w, Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    users,
	})

}

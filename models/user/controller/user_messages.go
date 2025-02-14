package controller

import "Petstore/models/user/entity"

type LoginRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool `json:"success" example:"true"`
	Code    int  `json:"code,omitempty" example:"200"`
	Data    Data `json:"data" example:"true"`
}

type Data struct {
	Message string      `json:"message" example:"success"`
	User    entity.User `json:"user" example:"true"`
}

type Response struct {
	Success bool        `json:"success" example:"true"`
	Code    int         `json:"code,omitempty" example:"200"`
	Data    interface{} `json:"data" example:"data"`
}

type LogoutResponse struct {
	Success bool   `json:"success" example:"true"`
	Code    int    `json:"code,omitempty" example:"200"`
	Data    string `json:"data" example:"Logiut success"`
}

type GetUserResponse struct {
	Success bool        `json:"success" example:"true"`
	Code    int         `json:"code,omitempty" example:"200"`
	User    entity.User `json:"user"`
}

type GetUsersResponse struct {
	Success bool          `json:"success" example:"true"`
	Code    int           `json:"code,omitempty" example:"200"`
	User    []entity.User `json:"user"`
}

package handlers

import (
	"fmt"
	"letterboxd-cineville/service"
	"letterboxd-cineville/views"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *service.Service
}

func NewUserHandler(service *service.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) HandleGetUsers(c echo.Context) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching users")
	}

	return views.Render(c, views.UserList(users))
}

func (h *UserHandler) HandleCreateUser(c echo.Context) error {
	email := c.FormValue("email")
	username := c.FormValue("username")

	err := h.service.CreateNewUser(email, username)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Fetch the updated list of users from the database
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching users")
	}

	// Render the updated user list only
	return views.Render(c, views.UserListOnly(users))
}

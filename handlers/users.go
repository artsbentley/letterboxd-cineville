package handlers

import (
	"letterboxd-cineville/db"
	"letterboxd-cineville/model"
	"letterboxd-cineville/views"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	db *db.Store
}

func NewUserHandler(db *db.Store) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) HandleGetUsers(c echo.Context) error {
	users, err := h.db.GetAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching users")
	}

	return views.Render(c, views.UserList(users))
}

func (h *UserHandler) HandleCreateUser(c echo.Context) error {
	email := c.FormValue("email")
	username := c.FormValue("username")

	user := model.User{
		Email:              email,
		LetterboxdUsername: username,
		Watchlist:          make([]string, 0),
	}

	// Insert the user into the database
	err := h.db.InsertWatchlist(user)
	if err != nil {
		// Handle duplicate error case and return appropriate response
		if strings.Contains(err.Error(), "already exists") {
			return c.String(http.StatusConflict, err.Error()) // Return 409 Conflict for duplicates
		}
		return c.String(http.StatusInternalServerError, "Error creating user")
	}

	// Fetch the updated list of users from the database
	users, err := h.db.GetAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching users")
	}

	// Render the updated user list only
	return views.Render(c, views.UserListOnly(users))
}

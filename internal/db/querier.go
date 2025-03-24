// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AssignUserLocation(ctx context.Context, arg AssignUserLocationParams) error
	CreateFilmEvent(ctx context.Context, arg CreateFilmEventParams) (FilmEvent, error)
	CreateLocation(ctx context.Context, city string) (uuid.UUID, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteFilmEvent(ctx context.Context, id uuid.UUID) error
	DeletePastFilmEvents(ctx context.Context) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	DeleteUserWatchlist(ctx context.Context, email string) error
	GetFilmEventByID(ctx context.Context, id uuid.UUID) (FilmEvent, error)
	GetFilmEventsByUserEmail(ctx context.Context, email string) ([]FilmEvent, error)
	GetLocationByCity(ctx context.Context, city string) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserLocationCities(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetUserWatchlist(ctx context.Context, email string) ([]string, error)
	GetUsers(ctx context.Context) ([]User, error)
	ListFilmEvents(ctx context.Context) ([]FilmEvent, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserWatchlist(ctx context.Context, arg UpdateUserWatchlistParams) error
}

var _ Querier = (*Queries)(nil)

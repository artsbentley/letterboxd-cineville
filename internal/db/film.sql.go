// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: film.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFilmEvent = `-- name: CreateFilmEvent :one
INSERT INTO film_event (
    name, url, start_date, end_date, location_name, location_address, organizer_name, organizer_url, performer_name
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, name, url, start_date, end_date, location_name, location_address, organizer_name, organizer_url, performer_name
`

type CreateFilmEventParams struct {
	Name            string    `json:"name"`
	Url             string    `json:"url"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	LocationName    string    `json:"location_name"`
	LocationAddress string    `json:"location_address"`
	OrganizerName   string    `json:"organizer_name"`
	OrganizerUrl    string    `json:"organizer_url"`
	PerformerName   string    `json:"performer_name"`
}

func (q *Queries) CreateFilmEvent(ctx context.Context, arg CreateFilmEventParams) (FilmEvent, error) {
	row := q.db.QueryRow(ctx, createFilmEvent,
		arg.Name,
		arg.Url,
		arg.StartDate,
		arg.EndDate,
		arg.LocationName,
		arg.LocationAddress,
		arg.OrganizerName,
		arg.OrganizerUrl,
		arg.PerformerName,
	)
	var i FilmEvent
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.StartDate,
		&i.EndDate,
		&i.LocationName,
		&i.LocationAddress,
		&i.OrganizerName,
		&i.OrganizerUrl,
		&i.PerformerName,
	)
	return i, err
}

const deleteFilmEvent = `-- name: DeleteFilmEvent :exec
DELETE FROM film_event
WHERE id = $1
`

func (q *Queries) DeleteFilmEvent(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteFilmEvent, id)
	return err
}

const deletePastFilmEvents = `-- name: DeletePastFilmEvents :exec
DELETE FROM film_event
WHERE start_date < NOW()
`

func (q *Queries) DeletePastFilmEvents(ctx context.Context) error {
	_, err := q.db.Exec(ctx, deletePastFilmEvents)
	return err
}

const getFilmEventByID = `-- name: GetFilmEventByID :one
SELECT id, name, url, start_date, end_date, location_name, location_address, organizer_name, organizer_url, performer_name
FROM film_event
WHERE id = $1
`

func (q *Queries) GetFilmEventByID(ctx context.Context, id uuid.UUID) (FilmEvent, error) {
	row := q.db.QueryRow(ctx, getFilmEventByID, id)
	var i FilmEvent
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.StartDate,
		&i.EndDate,
		&i.LocationName,
		&i.LocationAddress,
		&i.OrganizerName,
		&i.OrganizerUrl,
		&i.PerformerName,
	)
	return i, err
}

const getFilmEventsByUserEmail = `-- name: GetFilmEventsByUserEmail :many
SELECT fe.id, fe.name, fe.url, fe.start_date, fe.end_date, fe.location_name, fe.location_address, fe.organizer_name, fe.organizer_url, fe.performer_name
FROM film_event fe
JOIN users u ON fe.name = ANY(u.watchlist)
WHERE u.email = $1
`

func (q *Queries) GetFilmEventsByUserEmail(ctx context.Context, email string) ([]FilmEvent, error) {
	rows, err := q.db.Query(ctx, getFilmEventsByUserEmail, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FilmEvent{}
	for rows.Next() {
		var i FilmEvent
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.StartDate,
			&i.EndDate,
			&i.LocationName,
			&i.LocationAddress,
			&i.OrganizerName,
			&i.OrganizerUrl,
			&i.PerformerName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFilmEvents = `-- name: ListFilmEvents :many
SELECT id, name, url, start_date, end_date, location_name, location_address, organizer_name, organizer_url, performer_name
FROM film_event
`

func (q *Queries) ListFilmEvents(ctx context.Context) ([]FilmEvent, error) {
	rows, err := q.db.Query(ctx, listFilmEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FilmEvent{}
	for rows.Next() {
		var i FilmEvent
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.StartDate,
			&i.EndDate,
			&i.LocationName,
			&i.LocationAddress,
			&i.OrganizerName,
			&i.OrganizerUrl,
			&i.PerformerName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

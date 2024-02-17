// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package repo

import (
	"time"
)

type Org struct {
	ID        int32     `json:"id"`
	UniqID    string    `json:"uniq_id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Radar struct {
	ID        int32     `json:"id"`
	UniqID    string    `json:"uniq_id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RadarItem struct {
	ID          int32     `json:"id"`
	UniqID      string    `json:"uniq_id"`
	RadarID     int32     `json:"radar_id"`
	QuadrantID  int32     `json:"quadrant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type RadarQuadrant struct {
	ID      int32  `json:"id"`
	UniqID  string `json:"uniq_id"`
	RadarID int32  `json:"radar_id"`
	Name    string `json:"name"`
}

type User struct {
	ID             int32     `json:"id"`
	UniqID         string    `json:"uniq_id"`
	OrganisationID int32     `json:"organisation_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

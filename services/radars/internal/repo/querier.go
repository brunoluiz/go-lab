// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package repo

import (
	"context"
)

type Querier interface {
	DeleteRadar(ctx context.Context, uniqID string) error
	DeleteRadarItem(ctx context.Context, uniqID string) error
	DeleteRadarQuadrantByID(ctx context.Context, radarID int32) error
	GetRadarByID(ctx context.Context, uniqID string) (Radar, error)
	GetRadarItemsByRadarID(ctx context.Context, radarID int32) ([]GetRadarItemsByRadarIDRow, error)
	GetRadars(ctx context.Context) ([]GetRadarsRow, error)
	SaveRadar(ctx context.Context, arg SaveRadarParams) (Radar, error)
	SaveRadarItem(ctx context.Context, arg SaveRadarItemParams) (RadarItem, error)
	SaveRadarQuadrant(ctx context.Context, arg SaveRadarQuadrantParams) (RadarQuadrant, error)
}

var _ Querier = (*Queries)(nil)

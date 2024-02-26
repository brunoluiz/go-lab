package handler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
)

func (h *Handler) GetRadarById(ctx context.Context, req openapi.GetRadarByIdRequestObject) (openapi.GetRadarByIdResponseObject, error) {
	r, err := h.Repo.GetRadarByID(ctx, req.RadarId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &app.ErrNotFound{}
		}
		return nil, err
	}

	out := &openapi.Radar{
		Id:        r.UniqID,
		Title:     r.Title,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Items:     []openapi.RadarItem{},
		Quadrants: []openapi.RadarQuadrant{},
	}

	ris, err := h.Repo.GetRadarItemsByRadarID(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	for _, ri := range ris {
		out.Items = append(out.Items, openapi.RadarItem{
			CreatedAt:   ri.RadarItem.CreatedAt,
			Description: ri.RadarItem.Description,
			Name:        ri.RadarItem.Name,
			Id:          ri.RadarItem.UniqID,
			UpdatedAt:   ri.RadarItem.UpdatedAt,
			Quadrant: openapi.RadarQuadrant{
				Name: ri.RadarQuadrant.Name,
				Id:   ri.RadarQuadrant.UniqID,
			},
		})
	}

	rqs, err := h.Repo.GetRadarQuadrantsByRadarID(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	for _, rq := range rqs {
		out.Quadrants = append(out.Quadrants, openapi.RadarQuadrant{
			Id:   rq.UniqID,
			Name: rq.Name,
		})
	}

	return openapi.GetRadarById200JSONResponse{
		Status: openapi.StatusSuccess,
		Data: &openapi.DataResponse{
			Radar: out,
		},
	}, nil
}

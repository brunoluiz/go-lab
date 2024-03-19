package handler

import (
	"context"

	"github.com/brunoluiz/go-lab/core/genid"
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
)

func (h *Handler) AddRadarItem(ctx context.Context, req openapi.AddRadarItemRequestObject) (openapi.AddRadarItemResponseObject, error) {
	r, err := h.Repo.GetRadarByID(ctx, req.RadarId)
	if err != nil {
		return nil, err
	}

	rq, err := h.Repo.GetRadarQuadrantByUniqID(ctx, req.Body.QuadrantId)
	if err != nil {
		return nil, err
	}

	ri, err := h.Repo.SaveRadarItem(ctx, repo.SaveRadarItemParams{
		UniqID:      genid.New(genid.EntityRadarItem),
		RadarID:     r.ID,
		QuadrantID:  rq.ID,
		Name:        req.Body.Name,
		Description: req.Body.Description,
	})
	if err != nil {
		return nil, err
	}

	return openapi.AddRadarItem201JSONResponse{
		Status: openapi.StatusSuccess,
		Data: &openapi.DataResponse{
			RadarItem: &openapi.RadarItem{
				Id:          ri.UniqID,
				Name:        ri.Name,
				Description: ri.Description,
				UpdatedAt:   ri.UpdatedAt,
				CreatedAt:   ri.CreatedAt,
				Quadrant: openapi.RadarQuadrant{
					Id:   rq.UniqID,
					Name: rq.Name,
				},
			},
		},
	}, nil
}

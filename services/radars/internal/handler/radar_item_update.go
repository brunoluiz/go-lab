package handler

import (
	"context"

	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
)

func (h *Handler) UpdateRadarItem(ctx context.Context, req openapi.UpdateRadarItemRequestObject) (openapi.UpdateRadarItemResponseObject, error) {
	r, err := h.Repo.GetRadarByID(ctx, req.RadarId)
	if err != nil {
		return nil, err
	}

	rq, err := h.Repo.GetRadarQuadrantByUniqID(ctx, req.Body.QuadrantId)
	if err != nil {
		return nil, err
	}

	ri, err := h.Repo.SaveRadarItem(ctx, repo.SaveRadarItemParams{
		UniqID:      req.RadarItemId,
		RadarID:     r.ID,
		QuadrantID:  rq.ID,
		Name:        req.Body.Name,
		Description: req.Body.Description,
	})
	if err != nil {
		return nil, err
	}

	return openapi.UpdateRadarItem200JSONResponse{
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

package handler

import (
	"context"

	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
)

func (h *Handler) DeleteRadarItem(ctx context.Context, req openapi.DeleteRadarItemRequestObject) (openapi.DeleteRadarItemResponseObject, error) {
	err := h.Repo.DeleteRadarItem(ctx, req.RadarItemId)
	if err != nil {
		return nil, err
	}

	return openapi.DeleteRadarItem200JSONResponse{
		Status: openapi.StatusSuccess,
	}, nil
}

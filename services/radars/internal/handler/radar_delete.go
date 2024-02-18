package handler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
)

func (h *Handler) DeleteRadar(ctx context.Context, req openapi.DeleteRadarRequestObject) (openapi.DeleteRadarResponseObject, error) {
	err := h.Repo.DeleteRadar(ctx, req.RadarId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &app.ErrNotFound{}
		}
		return nil, err
	}

	return openapi.DeleteRadar200JSONResponse{
		Success: true,
	}, nil
}

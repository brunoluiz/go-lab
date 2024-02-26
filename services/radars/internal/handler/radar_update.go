package handler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
)

func (h *Handler) UpdateRadar(ctx context.Context, req openapi.UpdateRadarRequestObject) (openapi.UpdateRadarResponseObject, error) {
	out, err := h.Repo.SaveRadar(ctx, repo.SaveRadarParams{
		UniqID: req.RadarId,
		Title:  req.Body.Title,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &app.ErrNotFound{}
		}
		return nil, err
	}

	return openapi.UpdateRadar200JSONResponse{
		Status: StatusSuccess,
		Data: &openapi.DataResponse{
			Radar: &openapi.Radar{
				Title:     out.Title,
				Id:        out.UniqID,
				CreatedAt: out.CreatedAt,
				UpdatedAt: out.UpdatedAt,
			},
		},
	}, nil
}

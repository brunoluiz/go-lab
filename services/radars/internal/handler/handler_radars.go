package handler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/core/genid"
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
)

func (h *Handler) AddRadar(ctx context.Context, req openapi.AddRadarRequestObject) (openapi.AddRadarResponseObject, error) {
	var res openapi.AddRadar201JSONResponse

	err := h.WithTx(func(tx repo.Querier) error {
		radar, err := tx.SaveRadar(ctx, repo.SaveRadarParams{
			UniqID: genid.New(genid.EntityRadar),
			Title:  req.Body.Title,
		})
		if err != nil {
			return err
		}
		res = openapi.AddRadar201JSONResponse{
			Title:     radar.Title,
			UniqId:    radar.UniqID,
			CreatedAt: radar.CreatedAt,
			UpdatedAt: radar.UpdatedAt,
		}

		quadrantParams := []repo.SaveRadarQuadrantParams{
			{UniqID: genid.New(genid.EntityRadarItem), Name: "Techniques"},
			{UniqID: genid.New(genid.EntityRadarItem), Name: "Platforms"},
			{UniqID: genid.New(genid.EntityRadarItem), Name: "Tools"},
			{UniqID: genid.New(genid.EntityRadarItem), Name: "Languages & Frameworks"},
		}
		for _, params := range quadrantParams {
			params.RadarID = radar.ID

			_, err := tx.SaveRadarQuadrant(ctx, params)
			if err != nil {
				return err
			}
			// TODO: assign items to response
		}

		return nil
	})
	if err != nil {
		return openapi.AddRadar201JSONResponse{}, err
	}

	return res, nil
}

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
		Title:     out.Title,
		UniqId:    out.UniqID,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}, nil
}

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

func (h *Handler) GetRadarById(ctx context.Context, req openapi.GetRadarByIdRequestObject) (openapi.GetRadarByIdResponseObject, error) {
	out, err := h.Repo.GetRadarByID(ctx, req.RadarId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &app.ErrNotFound{}
		}
		return nil, err
	}

	return openapi.GetRadarById200JSONResponse{
		UniqId:    out.UniqID,
		Title:     out.Title,
		CreatedAt: out.CreatedAt,
	}, nil
}

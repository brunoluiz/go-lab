package handler

import (
	"context"

	"github.com/brunoluiz/go-lab/core/genid"
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
)

func (h *Handler) AddRadar(ctx context.Context, req openapi.AddRadarRequestObject) (openapi.AddRadarResponseObject, error) {
	var res openapi.AddRadar201JSONResponse

	err := h.Tx.Run(func(tx repo.Querier) error {
		radar, err := tx.SaveRadar(ctx, repo.SaveRadarParams{
			UniqID: genid.New(genid.EntityRadar),
			Title:  req.Body.Title,
		})
		if err != nil {
			return err
		}
		res = openapi.AddRadar201JSONResponse{
			Title:     radar.Title,
			Id:        radar.UniqID,
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

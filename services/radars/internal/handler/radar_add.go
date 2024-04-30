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
		r, err := tx.SaveRadar(ctx, repo.SaveRadarParams{
			UniqID: genid.New(EntityRadar),
			Title:  req.Body.Title,
		})
		if err != nil {
			return err
		}

		radar := &openapi.Radar{
			Title:     r.Title,
			Id:        r.UniqID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			Items:     []openapi.RadarItem{},
			Quadrants: []openapi.RadarQuadrant{},
		}

		res = openapi.AddRadar201JSONResponse{
			Status: openapi.StatusSuccess,
			Data: &openapi.DataResponse{
				Radar: radar,
			},
		}

		quadrantParams := []repo.SaveRadarQuadrantParams{
			{UniqID: genid.New(EntityRadarItem), Name: "Techniques"},
			{UniqID: genid.New(EntityRadarItem), Name: "Platforms"},
			{UniqID: genid.New(EntityRadarItem), Name: "Tools"},
			{UniqID: genid.New(EntityRadarItem), Name: "Languages & Frameworks"},
		}
		for _, params := range quadrantParams {
			params.RadarID = r.ID

			rq, err := tx.SaveRadarQuadrant(ctx, params)
			if err != nil {
				return err
			}

			radar.Quadrants = append(radar.Quadrants, openapi.RadarQuadrant{
				Name: rq.Name,
				Id:   rq.UniqID,
			})
		}

		return nil
	})
	if err != nil {
		return openapi.AddRadar201JSONResponse{}, err
	}

	return res, nil
}

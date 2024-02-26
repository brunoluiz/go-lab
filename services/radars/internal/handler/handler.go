package handler

import (
	"github.com/brunoluiz/go-lab/services/radars/gen/openapi"
	"github.com/brunoluiz/go-lab/services/radars/internal/repo"
)

type Handler struct {
	openapi.StrictServerInterface
	Repo repo.Querier
	Tx   repo.Tx
}

func New(
	q repo.Querier,
	withTx repo.Tx,
) *Handler {
	return &Handler{
		Repo: q,
		Tx:   withTx,
	}
}

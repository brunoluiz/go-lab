package repo

import (
	"context"
	"time"

	"github.com/brunoluiz/go-lab/services/apollo"
)

type list struct {
}

func List() *list {
	return &list{}
}

func (l *list) ByID(ctx context.Context, id string) (*apollo.List, error) {
	return &apollo.List{
		ID:        "123",
		Title:     "Foo",
		CreatedAt: time.Now(),
	}, nil
}

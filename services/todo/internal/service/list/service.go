package list

import (
	"context"
	"log/slog"
	"time"

	"github.com/brunoluiz/go-lab/lib/errx"
	"github.com/brunoluiz/go-lab/services/todo/internal/model"
	"github.com/brunoluiz/go-lab/services/todo/internal/repo"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func toDtoList(l model.List) List {
	return List{
		ID:        l.ID,
		Name:      l.Name,
		CreatedAt: l.CreatedAt,
	}
}

func fromDtoList(l List) model.List {
	return model.List{
		ID:        l.ID,
		Name:      l.Name,
		CreatedAt: l.CreatedAt,
	}
}

type Service struct {
	listRepo  repo.ListRepository
	logger    *slog.Logger
	validator *validator.Validate
}

func NewService(listRepo repo.ListRepository, logger *slog.Logger, validator *validator.Validate) *Service {
	return &Service{
		listRepo:  listRepo,
		logger:    logger,
		validator: validator,
	}
}

func (s *Service) CreateList(ctx context.Context, req CreateListRequest) (CreateListResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return CreateListResponse{}, errx.ErrValidation.Wrap(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateListResponse{}, errx.ErrUnknown.Wrap(err)
	}
	list := model.List{
		ID:        id.String(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}
	created, err := s.listRepo.CreateList(ctx, list)
	if err != nil {
		return CreateListResponse{}, err
	}
	return CreateListResponse{List: toDtoList(created)}, nil
}

func (s *Service) GetList(ctx context.Context, req GetListRequest) (GetListResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return GetListResponse{}, errx.ErrValidation.Wrap(err)
	}

	list, err := s.listRepo.GetList(ctx, req.ListID)
	if err != nil {
		return GetListResponse{}, err
	}
	return GetListResponse{List: toDtoList(list)}, nil
}

func (s *Service) ListLists(ctx context.Context, _ ListListsRequest) (ListListsResponse, error) {
	lists, err := s.listRepo.ListLists(ctx)
	if err != nil {
		return ListListsResponse{}, err
	}
	dtoLists := make([]List, len(lists))
	for i, l := range lists {
		dtoLists[i] = toDtoList(l)
	}
	return ListListsResponse{Lists: dtoLists}, nil
}

func (s *Service) UpdateList(ctx context.Context, req UpdateListRequest) (UpdateListResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return UpdateListResponse{}, errx.ErrValidation.Wrap(err)
	}

	list := fromDtoList(req.List)
	updated, err := s.listRepo.UpdateList(ctx, list)
	if err != nil {
		return UpdateListResponse{}, err
	}
	return UpdateListResponse{List: toDtoList(updated)}, nil
}

func (s *Service) DeleteList(ctx context.Context, req DeleteListRequest) (DeleteListResponse, error) {
	if err := s.validator.StructCtx(ctx, req); err != nil {
		return DeleteListResponse{}, errx.ErrValidation.Wrap(err)
	}

	err := s.listRepo.DeleteList(ctx, req.ListID)
	if err != nil {
		return DeleteListResponse{}, err
	}
	return DeleteListResponse{}, nil
}

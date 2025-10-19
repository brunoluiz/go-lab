package list

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/brunoluiz/go-lab/lib/app"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/model"
	"github.com/brunoluiz/go-lab/services/todo/internal/database/repository"
	"github.com/brunoluiz/go-lab/services/todo/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var ErrNameRequired = errors.New("name is required")

func toDtoList(l model.List) dto.List {
	return dto.List{
		ID:        l.ID,
		Name:      l.Name,
		CreatedAt: l.CreatedAt,
	}
}

func fromDtoList(l dto.List) model.List {
	return model.List{
		ID:        l.ID,
		Name:      l.Name,
		CreatedAt: l.CreatedAt,
	}
}

type Service struct {
	listRepo  repository.ListRepository
	logger    *slog.Logger
	validator *validator.Validate
}

func NewService(listRepo repository.ListRepository, logger *slog.Logger) *Service {
	return &Service{
		listRepo:  listRepo,
		logger:    logger,
		validator: validator.New(),
	}
}

func (s *Service) CreateList(ctx context.Context, req dto.CreateListRequest) (dto.CreateListResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.CreateListResponse{}, fmt.Errorf("%w: %w", app.ErrValidation, err)
	}
	id, err := uuid.NewV7()
	if err != nil {
		return dto.CreateListResponse{}, fmt.Errorf("%w: %w", app.ErrUnknown, err)
	}
	list := model.List{
		ID:        id.String(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}
	created, err := s.listRepo.CreateList(ctx, list)
	if err != nil {
		return dto.CreateListResponse{}, err
	}
	return dto.CreateListResponse{List: toDtoList(created)}, nil
}

func (s *Service) GetList(ctx context.Context, req dto.GetListRequest) (dto.GetListResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.GetListResponse{}, fmt.Errorf("%w: %w", app.ErrValidation, err)
	}
	list, err := s.listRepo.GetList(ctx, req.ListID)
	if err != nil {
		return dto.GetListResponse{}, err
	}
	return dto.GetListResponse{List: toDtoList(list)}, nil
}

func (s *Service) ListLists(ctx context.Context, _ dto.ListListsRequest) (dto.ListListsResponse, error) {
	lists, err := s.listRepo.ListLists(ctx)
	if err != nil {
		return dto.ListListsResponse{}, err
	}
	dtoLists := make([]dto.List, len(lists))
	for i, l := range lists {
		dtoLists[i] = toDtoList(l)
	}
	return dto.ListListsResponse{Lists: dtoLists}, nil
}

func (s *Service) UpdateList(ctx context.Context, req dto.UpdateListRequest) (dto.UpdateListResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.UpdateListResponse{}, fmt.Errorf("%w: %w", app.ErrValidation, err)
	}
	list := fromDtoList(req.List)
	updated, err := s.listRepo.UpdateList(ctx, list)
	if err != nil {
		return dto.UpdateListResponse{}, err
	}
	return dto.UpdateListResponse{List: toDtoList(updated)}, nil
}

func (s *Service) DeleteList(ctx context.Context, req dto.DeleteListRequest) (dto.DeleteListResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return dto.DeleteListResponse{}, fmt.Errorf("%w: %w", app.ErrValidation, err)
	}
	err := s.listRepo.DeleteList(ctx, req.ListID)
	if err != nil {
		return dto.DeleteListResponse{}, err
	}
	return dto.DeleteListResponse{}, nil
}

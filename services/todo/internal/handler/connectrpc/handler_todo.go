package connectrpc

import (
	"context"

	v1 "github.com/brunoluiz/go-lab/gen/go/proto/acme/api/todo/v1"
	"github.com/brunoluiz/go-lab/services/todo/internal/service/list"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateList(ctx context.Context, req *v1.CreateListRequest) (*v1.CreateListResponse, error) {
	internalReq := list.CreateListRequest{Name: req.Name}
	resp, err := h.listService.CreateList(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	return &v1.CreateListResponse{List: toProtoList(resp.List)}, nil
}

func (h *Handler) GetList(ctx context.Context, req *v1.GetListRequest) (*v1.GetListResponse, error) {
	internalReq := list.GetListRequest{ListID: req.ListId}
	resp, err := h.listService.GetList(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	return &v1.GetListResponse{List: toProtoList(resp.List)}, nil
}

func (h *Handler) ListLists(ctx context.Context, _ *v1.ListListsRequest) (*v1.ListListsResponse, error) {
	internalReq := list.ListListsRequest{}
	resp, err := h.listService.ListLists(ctx, internalReq)
	if err != nil {
		return nil, err
	}

	protoLists := make([]*v1.List, len(resp.Lists))
	for i, l := range resp.Lists {
		protoLists[i] = toProtoList(l)
	}
	return &v1.ListListsResponse{Lists: protoLists}, nil
}

func (h *Handler) UpdateList(ctx context.Context, req *v1.UpdateListRequest) (*v1.UpdateListResponse, error) {
	internalReq := list.UpdateListRequest{List: fromProtoList(req.List)}
	resp, err := h.listService.UpdateList(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateListResponse{List: toProtoList(resp.List)}, nil
}

func (h *Handler) DeleteList(ctx context.Context, req *v1.DeleteListRequest) (*v1.DeleteListResponse, error) {
	internalReq := list.DeleteListRequest{ListID: req.ListId}
	_, err := h.listService.DeleteList(ctx, internalReq)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteListResponse{}, nil
}

func toProtoList(l list.List) *v1.List {
	return &v1.List{
		Id:        l.ID,
		Name:      l.Name,
		CreatedAt: timestamppb.New(l.CreatedAt),
	}
}

func fromProtoList(l *v1.List) list.List {
	return list.List{
		ID:        l.Id,
		Name:      l.Name,
		CreatedAt: l.CreatedAt.AsTime(),
	}
}

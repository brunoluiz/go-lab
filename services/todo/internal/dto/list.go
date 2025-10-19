package dto

import "time"

type CreateListRequest struct {
	Name string `validate:"required"`
}

type CreateListResponse struct {
	List List
}

type GetListRequest struct {
	ListID string `validate:"required"`
}

type GetListResponse struct {
	List List
}

type ListListsRequest struct {
	// Empty for now
}

type ListListsResponse struct {
	Lists []List
}

type UpdateListRequest struct {
	List List `validate:"required"`
}

type UpdateListResponse struct {
	List List
}

type DeleteListRequest struct {
	ListID string `validate:"required"`
}

type DeleteListResponse struct {
	// Empty
}

type List struct {
	ID        string    `validate:"required"`
	Name      string    `validate:"required"`
	CreatedAt time.Time `validate:"-"`
}

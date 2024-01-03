package pagination

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/leomirandadev/clean-architecture-go/pkg/customerr"
	"github.com/leomirandadev/clean-architecture-go/pkg/httprouter"
)

type RequestPagination struct {
	Page       int
	Size       int
	SearchTerm *string
}

type ResponsePagination[T any] struct {
	Results    []T `json:"results"`
	TotalPages int `json:"total_pages"`
}

const MAX_PER_PAGE = 50

func GetFromQueryParams(c httprouter.Context) (RequestPagination, error) {
	var (
		err        error
		pagination RequestPagination
	)

	pagination.Page, err = strconv.Atoi(c.GetQueryParam("page"))
	if err != nil {
		return pagination, customerr.WithStatus(http.StatusBadRequest, "error to parse page", err)
	}

	if pagination.Page < 0 {
		return pagination, customerr.WithStatus(http.StatusBadRequest, "page should be bigger than 0", nil)
	}

	pagination.Size, err = strconv.Atoi(c.GetQueryParam("size"))
	if err != nil {
		return pagination, customerr.WithStatus(http.StatusBadRequest, "error to parse size", err)
	}

	if pagination.Size < 0 {
		return pagination, customerr.WithStatus(http.StatusBadRequest, "size should be bigger than 0", nil)
	}

	if pagination.Size > MAX_PER_PAGE {
		return pagination, customerr.WithStatus(http.StatusBadRequest, fmt.Sprintf("size should be less than %v", MAX_PER_PAGE), nil)
	}

	searchTerm := c.GetQueryParam("search_term")
	if searchTerm != "" {
		pagination.SearchTerm = &searchTerm
	}

	return pagination, nil
}

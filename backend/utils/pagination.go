package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Sort     string `json:"sort"`
	Order    string `json:"order"` // asc or desc
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination contains pagination metadata
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// Default pagination values
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
	DefaultSort     = "created_at"
	DefaultOrder    = "desc"
)

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(c *gin.Context) PaginationParams {
	page := DefaultPage
	pageSize := DefaultPageSize
	sort := c.DefaultQuery("sort", DefaultSort)
	order := c.DefaultQuery("order", DefaultOrder)

	// Parse page
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}

	// Parse page_size
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = ps
		if pageSize > MaxPageSize {
			pageSize = MaxPageSize
		}
	}

	// Validate order
	if order != "asc" && order != "desc" {
		order = DefaultOrder
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	}
}

// GetOffset calculates the database offset
func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the page size (limit)
func (p *PaginationParams) GetLimit() int {
	return p.PageSize
}

// GetOrderClause returns the SQL ORDER BY clause
func (p *PaginationParams) GetOrderClause() string {
	return p.Sort + " " + p.Order
}

// BuildPagination creates pagination metadata from params and total count
func BuildPagination(params PaginationParams, totalItems int64) Pagination {
	totalPages := int((totalItems + int64(params.PageSize) - 1) / int64(params.PageSize))
	if totalPages < 1 {
		totalPages = 1
	}

	return Pagination{
		CurrentPage: params.Page,
		PageSize:    params.PageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     params.Page < totalPages,
		HasPrev:     params.Page > 1,
	}
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, params PaginationParams, totalItems int64) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: BuildPagination(params, totalItems),
	}
}

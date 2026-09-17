package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type Params struct {
	Page   int
	Limit  int
	Search string
	SortBy string
	Order  string
	Offset int
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// Parse reads page/limit/search/sort_by/order from query params.
// allowedSort is a whitelist of sortable columns to prevent SQL injection via sort_by.
func Parse(c *gin.Context, allowedSort []string, defaultSort string) Params {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = DefaultPage
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	sortBy := c.DefaultQuery("sort_by", defaultSort)
	valid := false
	for _, s := range allowedSort {
		if s == sortBy {
			valid = true
			break
		}
	}
	if !valid {
		sortBy = defaultSort
	}

	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return Params{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
		SortBy: sortBy,
		Order:  order,
		Offset: (page - 1) * limit,
	}
}

func BuildMeta(page, limit int, totalItems int64) Meta {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	return Meta{Page: page, Limit: limit, TotalItems: totalItems, TotalPages: totalPages}
}

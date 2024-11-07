package data

import (
	"strings"

	"github.com/hayohtee/social/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string
	Search   string
	Tags     []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(strings.ToUpper(f.Sort) == "ASC" || strings.ToUpper(f.Sort) == "DESC", "sort", "invalid sort value")
	v.Check(len(f.Tags) >= 0 || len(f.Tags) <= 5, "tags", "must contain a maximum of five tags")
	v.Check(len(f.Search) <= 100, "search", "must not be 100 bytes long")
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

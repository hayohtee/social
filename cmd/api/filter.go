package main

import (
	"strings"

	"github.com/hayohtee/social/internal/validator"
)

type filters struct {
	Page         int
	PageSize     int
	Sort         string
}

func ValidateFilters(v *validator.Validator, f filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(strings.ToUpper(f.Sort) == "ASC" || strings.ToLower(f.Sort) == "DESC", "sort", "invalid sort value")
}

func (f filters) limit() int {
	return f.PageSize
}

func (f filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

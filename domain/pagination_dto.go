package domain

type PaginationDTO struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

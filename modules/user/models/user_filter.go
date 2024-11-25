package models

import "fiber-ngulik/pkg/utils"

type UserFilter struct {
	Role string `json:"role" query:"role"`
	utils.PaginationRequest
}

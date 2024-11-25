package utils

import (
	"encoding/json"
	"fiber-ngulik/pkg/httperror"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Result struct {
	Data  interface{}
	Error interface{}
	Total int64
}

// Pagination data structure
type PaginationRequest struct {
	Page              int64 `json:"page" query:"page"`
	PerPage           int64 `json:"per_page" query:"per_page"`
	DisablePagination bool  `json:"disable_pagination" query:"disable_pagination"`
}

type PaginationResponse struct {
	Total   int64 `json:"total"`
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
}

type BaseWrapperModel struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Meta    interface{} `json:"meta,omitempty"`
}

type Meta struct {
	Method        string    `json:"method"`
	Url           string    `json:"url"`
	Code          string    `json:"code"`
	ContentLength int       `json:"content_length"`
	Date          time.Time `json:"date"`
	Ip            string    `json:"ip"`
}

func (q *PaginationRequest) GetOffset() int64 {
	if q.Page <= 1 {
		return 0
	}
	return (q.Page - 1) * q.PerPage
}

func (q *PaginationRequest) GetLimit() int64 {
	return q.PerPage
}

func (q *PaginationRequest) SetDefault() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PerPage == 0 {
		q.PerPage = 10
	}
}

func (q *PaginationRequest) SetPaginationResponse(page int64, perPage int64) PaginationRequest {
	return PaginationRequest{
		Page:    page,
		PerPage: perPage,
	}
}

func (q *PaginationRequest) GetPaginationRequest() PaginationRequest {
	return PaginationRequest{
		Page:              q.Page,
		PerPage:           q.PerPage,
		DisablePagination: q.DisablePagination,
	}
}

func Response(data interface{}, message string, code int, c *fiber.Ctx) error {
	success := false
	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Method(),
		Code:          fmt.Sprintf("%v", http.StatusOK),
		ContentLength: c.Request().Header.ContentLength(),
		Ip:            c.IP(),
	}
	byteMeta, _ := json.Marshal(meta)
	LogDefault(string(byteMeta))

	if code < http.StatusBadRequest {
		success = true
	}

	result := BaseWrapperModel{
		Success: success,
		Data:    data,
		Message: message,
		Code:    code,
	}

	return c.Status(code).JSON(result)
}

func ResponseWithPagination(data interface{}, message string, code int, total int64, pagination PaginationRequest, c *fiber.Ctx) error {
	success := false
	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Method(),
		Code:          fmt.Sprintf("%v", http.StatusOK),
		ContentLength: c.Request().Header.ContentLength(),
		Ip:            c.IP(),
	}
	byteMeta, _ := json.Marshal(meta)
	LogDefault(string(byteMeta))

	if code < http.StatusBadRequest {
		success = true
	}

	result := BaseWrapperModel{
		Success: success,
		Data:    data,
		Message: message,
		Code:    code,
	}

	if !pagination.DisablePagination {
		result.Meta = PaginationResponse{
			Total:   total,
			Page:    pagination.Page,
			PerPage: pagination.PerPage,
		}
	}

	return c.Status(code).JSON(result)
}

func ResponseError(err interface{}, c *fiber.Ctx) error {

	errObj := getErrStatusCode(err)

	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Method(),
		Code:          fmt.Sprintf("%v", errObj),
		Ip:            c.IP(),
		ContentLength: c.Request().Header.ContentLength(),
	}

	result := BaseWrapperModel{
		Success: false,
		Message: errObj.Message,
		Code:    errObj.Code,
	}

	byteMeta, _ := json.Marshal(meta)

	LogError(string(byteMeta))

	return c.Status(errObj.Code).JSON(result)
}

func LogDefault(meta string) {
	log.Default().Println("service-info", "Logging service...", "audit-log", meta)
}

func LogError(meta string) {
	log.Default().Println("service-error", "Logging service...", "audit-log", meta)
}

func getErrStatusCode(err interface{}) httperror.CommonErrorData {
	errData := httperror.CommonErrorData{}

	if e, ok := err.(httperror.ErrorString); ok {
		errData.ResponseCode = e.Code() // HTTP status code
		errData.Code = e.Code()         // Custom application code (same as HTTP for now)
		errData.Message = e.Message()   // Error message
		return errData
	}

	// Default case for unknown error types
	errData.Code = http.StatusInternalServerError
	errData.Message = "Unknown Error"
	return errData
}

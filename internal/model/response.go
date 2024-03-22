package model

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var response = map[int]BaseResponse{
	http.StatusOK:                  {http.StatusOK, "000", "Successful"},
	http.StatusCreated:             {http.StatusCreated, "000", "Successful"},
	http.StatusInternalServerError: {http.StatusInternalServerError, "001", "Internal Server Error"},
	http.StatusBadRequest:          {http.StatusBadRequest, "001", "Bad Request"},
}

type BaseResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ResponsePage struct {
	TotalData int64 `json:"total_data"`
	TotalPage int64 `json:"total_page"`
}

type ResponseData struct {
	BaseResponse
	*ResponsePage `json:",omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

type ResponseDataList struct {
	BaseResponse
	Data interface{} `json:"data,omitempty"`
}

func NewResponse(status int, data ...interface{}) (r ResponseData) {
	res, ok := response[status]

	if !ok {
		res = response[http.StatusInternalServerError]
		return ResponseData{
			BaseResponse: res,
		}
	}

	res.Status = status
	if len(data) > 0 {
		if _, ok := data[0].([]interface{}); !ok {
			r.Data = data[0]
		}

		if len(data) > 1 {
			if p, ok := data[1].(ResponsePage); ok {
				r.ResponsePage = &p
			}
		}
	}

	r.BaseResponse = res
	return r
}

func Response(ctx *fiber.Ctx, status int, data ...interface{}) error {
	return ctx.Status(status).JSON(NewResponse(status, data...))
}

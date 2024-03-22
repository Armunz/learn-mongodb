package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/Armunz/learn-mongodb/internal/model"
	"github.com/Armunz/learn-mongodb/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type resource struct {
	service  services.Service
	validate *validator.Validate
	timeout  int
}

func RegisterHandlers(r fiber.Router, service services.Service, validate *validator.Validate, timeout int) {
	res := resource{
		service:  service,
		validate: validate,
		timeout:  timeout,
	}

	r.Post("/", res.Create)
	r.Get("/", res.Get)
	r.Get("/:id", res.Detail)
	r.Put("/:id", res.Update)
	r.Delete("/:id", res.Delete)
}

func (r *resource) Create(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	var request model.AccountCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	if err := r.validate.Struct(request); err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	if err := r.service.CreateAccount(c.UserContext(), request); err != nil {
		return model.Response(c, fiber.StatusInternalServerError)
	}

	return model.Response(c, fiber.StatusCreated, nil)
}

func (r *resource) Get(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	var request model.AccountListRequest
	if err := c.QueryParser(&request); err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	response, totalData, totalPage, err := r.service.GetListAccount(c.UserContext(), request)
	if err != nil {
		return model.Response(c, fiber.StatusInternalServerError)
	}

	responsePage := model.ResponsePage{
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return model.Response(c, fiber.StatusOK, response, responsePage)
}

func (r *resource) Detail(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	accountID := c.Params("id")
	accountIDNum, err := strconv.Atoi(accountID)
	if err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	response, err := r.service.GetAccountDetail(c.UserContext(), accountIDNum)
	if err != nil {
		return model.Response(c, fiber.StatusInternalServerError)
	}

	return model.Response(c, fiber.StatusOK, response)
}

func (r *resource) Update(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	accountID := c.Params("id")
	accountIDNum, err := strconv.Atoi(accountID)
	if err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	var request model.AccountUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	if err := r.validate.Struct(request); err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	if err := r.service.UpdateAccount(c.UserContext(), accountIDNum, request); err != nil {
		return model.Response(c, fiber.StatusInternalServerError)
	}

	return model.Response(c, fiber.StatusOK)
}

func (r *resource) Delete(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	accountID := c.Params("id")
	accountIDNum, err := strconv.Atoi(accountID)
	if err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	if err := r.service.DeleteAccount(c.UserContext(), accountIDNum); err != nil {
		return model.Response(c, fiber.StatusInternalServerError)
	}

	return model.Response(c, fiber.StatusOK)
}

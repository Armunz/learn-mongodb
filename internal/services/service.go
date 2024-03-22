package services

import (
	"context"
	"strings"

	"github.com/Armunz/learn-mongodb/internal/entity"
	"github.com/Armunz/learn-mongodb/internal/model"
	"github.com/Armunz/learn-mongodb/internal/repositories"
)

type Service interface {
	CreateAccount(ctx context.Context, request model.AccountCreateRequest) error
	GetListAccount(ctx context.Context, request model.AccountListRequest) ([]model.AccountResponse, int64, int64, error)
	GetAccountDetail(ctx context.Context, accountID int) (model.AccountResponse, error)
	UpdateAccount(ctx context.Context, accountID int, request model.AccountUpdateRequest) error
	DeleteAccount(ctx context.Context, accountID int) error
}

type serviceImpl struct {
	repo         repositories.Repository
	defaultLimit int
}

func NewService(repo repositories.Repository, defaultLimit int) Service {
	return &serviceImpl{
		repo:         repo,
		defaultLimit: defaultLimit,
	}
}

// CreateAccount implements Service.
func (s *serviceImpl) CreateAccount(ctx context.Context, request model.AccountCreateRequest) error {
	account := entity.Account{
		AccountID: request.AccountID,
		Limit:     request.Limit,
		Products:  request.Products,
	}

	return s.repo.Create(ctx, account)
}

// DeleteAccount implements Service.
func (s *serviceImpl) DeleteAccount(ctx context.Context, accountID int) error {
	return s.repo.Delete(ctx, accountID)
}

// GetAccountDetail implements Service.
func (s *serviceImpl) GetAccountDetail(ctx context.Context, accountID int) (model.AccountResponse, error) {
	account, err := s.repo.GetByAccountID(ctx, accountID)
	if err != nil {
		return model.AccountResponse{}, err
	}

	response := model.AccountResponse{
		AccountID: account.AccountID,
		Limit:     account.Limit,
		Products:  account.Products,
	}

	return response, nil
}

// GetListAccount implements Service.
func (s *serviceImpl) GetListAccount(ctx context.Context, request model.AccountListRequest) ([]model.AccountResponse, int64, int64, error) {
	limit := request.Limit
	if limit == 0 {
		limit = s.defaultLimit
	}

	var offset int
	if request.Page > 0 {
		offset = (request.Page - 1) * limit
	}

	orderBy, err := validateOrderByRequest(request.OrderBy)
	if err != nil {
		return nil, 0, 0, err
	}

	accounts, count, err := s.repo.List(ctx, request.Product, orderBy, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	// count total pages
	var totalPages int64
	if limit > 0 {
		totalPages = count / int64(limit)
		if count%int64(limit) != 0 {
			totalPages++
		}
	}

	response := make([]model.AccountResponse, len(accounts))
	for i, a := range accounts {
		response[i] = model.AccountResponse{
			AccountID: a.AccountID,
			Limit:     a.Limit,
			Products:  a.Products,
		}
	}

	return response, count, totalPages, nil
}

// UpdateAccount implements Service.
func (s *serviceImpl) UpdateAccount(ctx context.Context, accountID int, request model.AccountUpdateRequest) error {
	account, err := s.repo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}

	account.Limit = request.Limit
	account.Products = request.Products

	return s.repo.Update(ctx, account)
}

func validateOrderByRequest(orderBy model.OrderField) (int, error) {
	if orderBy.AccountID != "" {
		if strings.EqualFold(orderBy.AccountID, ORDER_BY_ASC) {
			return 1, nil
		}

		if strings.EqualFold(orderBy.AccountID, ORDER_BY_DESC) {
			return -1, nil
		}

		return 0, errOrderByInvalid
	}

	return 0, nil
}

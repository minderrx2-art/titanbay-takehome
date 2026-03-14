package service

import (
	"titanbay/internal/domain"
	"titanbay/internal/store/postgres"

	"github.com/google/uuid"
)

type FundManager struct {
	store *postgres.Store
}

func NewFundManager(store *postgres.Store) *FundManager {
	return &FundManager{
		store: store,
	}
}

func (f *FundManager) GetAllFunds() ([]domain.Fund, error) {
	return f.store.GetAllFunds()
}

func (f *FundManager) CreateFund(input domain.Fund) (domain.Fund, error) {
	return f.store.CreateFund(input)
}

func (f *FundManager) UpdateFund(id uuid.UUID, input domain.Fund) (domain.Fund, error) {
	return f.store.UpdateFund(id, input)
}

func (f *FundManager) GetFundByID(id uuid.UUID) (domain.Fund, error) {
	return f.store.GetFundByID(id)
}

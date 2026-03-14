package service

import (
	"titanbay/internal/domain"
	"titanbay/internal/store/postgres"
)

type InvestorManager struct {
	store *postgres.Store
}

func NewInvestorManager(store *postgres.Store) *InvestorManager {
	return &InvestorManager{
		store: store,
	}
}

func (i *InvestorManager) GetAllInvestors() ([]domain.Investor, error) {
	return i.store.GetAllInvestors()
}

func (i *InvestorManager) CreateInvestor(input domain.Investor) (domain.Investor, error) {
	return i.store.CreateInvestor(input)
}

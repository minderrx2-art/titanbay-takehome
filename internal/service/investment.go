package service

import (
	"titanbay/internal/domain"
	"titanbay/internal/store/postgres"

	"github.com/google/uuid"
)

type InvestmentManager struct {
	store *postgres.Store
}

func NewInvestmentManager(store *postgres.Store) *InvestmentManager {
	return &InvestmentManager{
		store: store,
	}
}

func (i *InvestmentManager) GetInvestmentsByFund(fundID uuid.UUID) ([]domain.Investment, error) {
	return i.store.GetInvestmentsByFund(fundID)
}

func (i *InvestmentManager) CreateInvestment(input domain.Investment) (domain.Investment, error) {
	return i.store.CreateInvestment(input)
}

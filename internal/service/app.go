package service

import "titanbay/internal/store/postgres"

type App struct {
	Funds       *FundManager
	Investors   *InvestorManager
	Investments *InvestmentManager
}

func NewService(store *postgres.Store) *App {
	return &App{
		Funds:       NewFundManager(store),
		Investors:   NewInvestorManager(store),
		Investments: NewInvestmentManager(store),
	}
}

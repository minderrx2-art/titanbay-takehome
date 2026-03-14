package postgres

import (
	"fmt"
	"titanbay/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Go receiver pattern
type Store struct {
	pgdb *gorm.DB
}

func NewStoreModel(db *gorm.DB) *Store {
	return &Store{pgdb: db}
}

func (s *Store) Close() error {
	pgDB, err := s.pgdb.DB()
	if err != nil {
		return fmt.Errorf("database failed to close")
	}

	return pgDB.Close()
}

// List all funds
func (s *Store) GetAllFunds() ([]domain.Fund, error) {
	var funds []domain.Fund
	result := s.pgdb.Find(&funds)

	if result.Error != nil {
		return nil, result.Error
	}

	return funds, nil
}

// Create a new fund
func (s *Store) CreateFund(input domain.Fund) (domain.Fund, error) {
	result := s.pgdb.Create(&input)

	if result.Error != nil {
		return domain.Fund{}, result.Error
	}

	return input, nil
}

// Update an existing fund
func (s *Store) UpdateFund(id uuid.UUID, updates domain.Fund) (domain.Fund, error) {
	var fund domain.Fund

	if err := s.pgdb.First(&fund, "id = ?", id).Error; err != nil {
		return domain.Fund{}, err
	}

	if err := s.pgdb.Model(&fund).Updates(updates).Error; err != nil {
		return domain.Fund{}, err
	}

	return fund, nil
}

// Get a specific fund
func (s *Store) GetFundByID(id uuid.UUID) (domain.Fund, error) {
	var fund domain.Fund

	result := s.pgdb.First(&fund, "id = ?", id)

	if result.Error != nil {
		return domain.Fund{}, result.Error
	}

	return fund, nil
}

// List all investors
func (s *Store) GetAllInvestors() ([]domain.Investor, error) {
	var investors []domain.Investor
	if err := s.pgdb.Find(&investors).Error; err != nil {
		return nil, err
	}
	return investors, nil
}

// Create a new investor
func (s *Store) CreateInvestor(input domain.Investor) (domain.Investor, error) {
	if err := s.pgdb.Create(&input).Error; err != nil {
		return domain.Investor{}, err
	}
	return input, nil
}

// List all investments for a specific fund
func (s *Store) GetInvestmentsByFund(fundID uuid.UUID) ([]domain.Investment, error) {
	var investments []domain.Investment

	result := s.pgdb.Where("fund_id = ?", fundID).Find(&investments)

	if result.Error != nil {
		return nil, result.Error
	}

	return investments, nil
}

// Create a new investment
func (s *Store) CreateInvestment(input domain.Investment) (domain.Investment, error) {
	result := s.pgdb.Create(&input)

	if result.Error != nil {
		return domain.Investment{}, result.Error
	}

	return input, nil
}

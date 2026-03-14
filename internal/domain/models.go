package domain

import (
	"time"

	"github.com/google/uuid"
)

type Fund struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name          string       `json:"name"`
	VintageYear   int          `json:"vintage_year"`
	TargetSizeUSD float64      `json:"target_size_usd"`
	Status        string       `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	Investments   []Investment `gorm:"foreignKey:FundID" json:"-"`
}

type Investor struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name         string       `json:"name"`
	InvestorType string       `json:"investor_type"`
	Email        string       `json:"email"`
	CreatedAt    time.Time    `json:"created_at"`
	Investments  []Investment `gorm:"foreignKey:InvestorID" json:"-"`
}

type Investment struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	InvestorID     uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"investor_id"`
	FundID         uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"fund_id"`
	AmountUSD      float64   `json:"amount_usd"`
	InvestmentDate string    `gorm:"type:date" json:"investment_date"`
}

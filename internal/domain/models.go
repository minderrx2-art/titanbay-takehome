package domain

import (
	"time"

	"github.com/google/uuid"
)

// ISSUE: Override Go float64 into DB numeric, issue it will round to nearest
type Fund struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name          *string      `gorm:"not null" json:"name"`
	VintageYear   *int         `gorm:"not null" json:"vintage_year"`
	TargetSizeUSD *float64     `gorm:"type:numeric(15,2);not null" json:"target_size_usd"`
	Status        *string      `gorm:"not null" json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	Investments   []Investment `gorm:"foreignKey:FundID" json:"-"`
}

type Investor struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name         *string      `gorm:"not null" json:"name"`
	InvestorType *string      `gorm:"not null" json:"investor_type"`
	Email        *string      `gorm:"not null" json:"email"`
	CreatedAt    time.Time    `json:"created_at"`
	Investments  []Investment `gorm:"foreignKey:InvestorID" json:"-"`
}

// ISSUE: Override Go float64 into DB numeric, issue it will round to nearest
type Investment struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	InvestorID     uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"investor_id"`
	FundID         uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"fund_id"`
	AmountUSD      *float64  `gorm:"type:numeric(15,2);not null" json:"amount_usd"`
	InvestmentDate *string   `gorm:"type:date;not null" json:"investment_date"`
}

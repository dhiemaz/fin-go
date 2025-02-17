package entities

import (
	"time"
)

type Account struct {
	ID         int64    `gorm:"type:bigint;primary_key;"`
	CIF        string   `gorm:"type:char(36);not null"`
	NickName   string   `json:"nick_name"`
	Amount     float64  `gorm:"default:0.0;not_null" json:"amount"`
	CustomerID int64    `gorm:"type:bigint;not_null" json:"customer_id"`
	Customer   Customer `json:"customer"`
	//Transactions   []Transaction `json:"transactions"`
	//ToTransactions []Transaction `gorm:"foreignkey:ToAccountID" json:"to_transactions"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

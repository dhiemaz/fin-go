package customer

import (
	"github.com/dhiemaz/fin-go/domain/account"
	"github.com/dhiemaz/fin-go/domain/transaction"
	"github.com/google/uuid"
	"time"
)

type CustomerModel struct {
	ID           uuid.UUID                      `gorm:"type:char(36);primary_key;"`
	FirstName    string                         `gorm:"not_null" json:"first_name"`
	LastName     string                         `gorm:"not_null" json:"last_name"`
	Email        string                         `gorm:"unique;not_null" json:"email"`
	Password     string                         `gorm:"size:100;not null;" json:"password"`
	Accounts     []account.AccountModel         `json:"accounts"`
	Transactions []transaction.TransactionModel `json:"transactions"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

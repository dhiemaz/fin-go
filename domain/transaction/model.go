package transaction

import (
	"errors"
	"fmt"
	"github.com/dhiemaz/fin-go/domain/account"
	"github.com/dhiemaz/fin-go/domain/customer"
	"github.com/google/uuid"
	"time"
)

type TransactionModel struct {
	ID              uuid.UUID              `gorm:"type:char(36);primary_key;"`
	Amount          float64                `gorm:"default:0.0;not_null" json:"amount"`
	TransactionType string                 `gorm:"not_null" json:"transaction_type"`
	Notes           string                 `json:"notes"`
	AccountID       uuid.UUID              `gorm:"not_null" json:"account_id"`
	Account         account.AccountModel   `json:"account"`
	ToAccountID     uuid.UUID              `json:"to_account_id"`
	ToAccount       account.AccountModel   `gorm:"foreignkey:ToAccountID" json:"to_account"`
	CustomerID      uuid.UUID              `json:"customer_id"`
	Customer        customer.CustomerModel `json:"customer"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (transaction *TransactionModel) BeforeCreate(scope *gorm.Scope) error {
	uuidV4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidV4)
}

func (transaction *TransactionModel) BeforeSave() (err error) {
	if transaction.TransactionType != "withdraw" && transaction.TransactionType != "transfer" && transaction.TransactionType != "deposit" {
		err = errors.New("can't save without a correct transaction_type")
	}

	if transaction.TransactionType == "transfer" && fmt.Sprint(transaction.ToAccountID) == "" {
		err = errors.New("can't save without a to_account_id for transaction_type 'transfer'")
	}
	return
}

func (transaction TransactionModel) AfterCreate(db *gorm.DB) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		if transaction.TransactionType == "withdraw" || transaction.TransactionType == "transfer" {
			var account Account
			db.Model(&transaction).Association("Account").Find(&account)
			account.Amount -= transaction.Amount
			db.Save(&account)
			if transaction.TransactionType == "transfer" {
				var toAccount Account
				db.Model(&transaction).Association("ToAccount").Find(&toAccount)
				toAccount.Amount += transaction.Amount
				db.Save(&toAccount)
			}
		} else {
			var account Account
			db.Model(&transaction).Association("Account").Find(&account)
			account.Amount += transaction.Amount
			db.Save(&account)
		}

		return nil
	})
}

package handlers

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type AccountHandler struct {
	logger *zap.Logger
}

func NewAccountController(db *gorm.DB) *AccountHandler {
	return &AccountHandler{
		//repository:  accountRepo.NewAccountRepository(db),
		//customerRepository: customerRepo.NewCustomerRepository(db),
		logger:             ,
	}
}

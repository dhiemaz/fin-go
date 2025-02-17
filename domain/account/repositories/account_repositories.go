package repositories

import "github.com/jinzhu/gorm"

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

//func (r *AccountRepository) CreateCustomer

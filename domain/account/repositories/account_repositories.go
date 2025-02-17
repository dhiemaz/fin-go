package repositories

import (
	"context"
	"errors"
	"github.com/dhiemaz/fin-go/entities"
	"github.com/jinzhu/gorm"
)

// AccountRepository interface
type AccountRepository interface {
	Create(ctx context.Context, account entities.Account) error
	Delete(ctx context.Context, account entities.Account) error
	GetAll(ctx context.Context, limit int, offset int) ([]entities.Account, error)
	GetByCIF(ctx context.Context, cif string) (entities.Account, error)
	GetDataById(ctx context.Context, customerId int64) (entities.Account, error)
	Count(ctx context.Context) (int64, error)
	ExistsRecord(ctx context.Context, field string, value string) (bool, error)
}

type Account struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *Account {
	return &Account{
		db: db,
	}
}

func (repo *Account) Create(ctx context.Context, account entities.Account) error {
	result := repo.db.Create(&account)
	return result.Error
}

func (repo *Account) Delete(ctx context.Context, account entities.Account) error {
	result := repo.db.Delete(&account)
	return result.Error
}

func (repo *Account) GetAll(ctx context.Context, limit int, offset int) ([]entities.Account, error) {
	var account []entities.Account
	result := repo.db.
		Table("accounts").
		Limit(limit).Offset(offset).
		Find(&account)
	return account, result.Error
}

func (repo *Account) GetByCIF(ctx context.Context, CIF string) (entities.Account, error) {
	var account entities.Account
	result := repo.db.Table("accounts").First(&account, CIF)
	if result.Error != nil {
		return entities.Account{}, result.Error
	}
	return account, result.Error
}

func (repo *Account) GetDataById(ctx context.Context, accountId int64) (entities.Account, error) {
	var account entities.Account
	result := repo.db.Table("accounts").First(&account, accountId)
	if result.Error != nil {
		return entities.Account{}, result.Error
	}
	return account, result.Error
}

func (repo *Account) Count(ctx context.Context) (int64, error) {
	var count int64
	result := repo.db.Table("accounts").Count(&count)
	return count, result.Error
}

// ExistsRecord : check if record exist by valid fields
func (repo *Account) ExistsRecord(ctx context.Context, field string, value string) (bool, error) {
	var count int64
	// Validate the field to avoid SQL injection
	validFields := map[string]bool{
		"cif":         true,
		"customer_id": true,
	}
	if !validFields[field] {
		return false, errors.New("invalid field name")
	}
	// Construct and execute the query
	query := repo.db.Table("accounts").Where(field+" = ?", value).Count(&count)
	if query.Error != nil {
		return false, query.Error
	}
	return count > 0, nil
}

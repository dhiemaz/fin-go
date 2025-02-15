package repositories

import (
	"context"
	"errors"
	"github.com/dhiemaz/fin-go/entities"
	"github.com/jinzhu/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (repo *CustomerRepository) Create(ctx context.Context, customer entities.Customer) error {
	result := repo.db.Create(&customer)
	return result.Error
}

func (repo *CustomerRepository) CreateBatch(ctx context.Context, customers []entities.Customer) error {
	tx := repo.db.Begin()
	result := repo.db.Create(&customers)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	return tx.Commit().Error
}

func (repo *CustomerRepository) Update(ctx context.Context, customer entities.Customer) error {
	result := repo.db.Save(&customer)
	return result.Error
}

func (repo *CustomerRepository) Delete(ctx context.Context, customer entities.Customer) error {
	result := repo.db.Delete(&customer)
	return result.Error
}

func (repo *CustomerRepository) GetAll(ctx context.Context, limit int, offset int) ([]entities.CustomerData, error) {
	var customers []entities.CustomerData
	result := repo.db.
		Table("view_customer_data").
		Limit(limit).Offset(offset).
		Find(&customers)
	return customers, result.Error
}

func (repo *CustomerRepository) GetById(ctx context.Context, customerId int64) (entities.Customer, error) {
	var customer entities.Customer
	result := repo.db.Table("customers").First(&customer, customerId)
	if result.Error != nil {
		return entities.Customer{}, result.Error
	}
	return customer, result.Error
}

func (repo *CustomerRepository) GetDataById(ctx context.Context, customerId int64) (entities.CustomerData, error) {
	var customer entities.CustomerData
	result := repo.db.Table("view_customer_data").First(&customer, customerId)
	if result.Error != nil {
		return entities.CustomerData{}, result.Error
	}
	return customer, result.Error
}

func (repo *CustomerRepository) GetByUniqueId(ctx context.Context, uniqueId string) (entities.Customer, error) {
	var customer entities.Customer
	result := repo.db.Table("customers").First(&customer, "unique_id", uniqueId)
	if result.Error != nil {
		return entities.Customer{}, result.Error
	}
	return customer, result.Error
}

func (repo *CustomerRepository) GetByDataUniqueId(ctx context.Context, uniqueId string) (entities.CustomerData, error) {
	var customer entities.CustomerData
	result := repo.db.Table("view_customer_data").First(&customer, "unique_id", uniqueId)
	if result.Error != nil {
		return entities.CustomerData{}, result.Error
	}
	return customer, result.Error
}

func (repo *CustomerRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	result := repo.db.Table("customers").Count(&count)
	return count, result.Error
}

func (repo *CustomerRepository) ExistsRecord(ctx context.Context, field string, value string) (bool, error) {
	var count int64
	// Validate the field to avoid SQL injection
	validFields := map[string]bool{
		"identification_number": true,
		"email":                 true,
		"phone":                 true,
	}
	if !validFields[field] {
		return false, errors.New("invalid field name")
	}
	// Construct and execute the query
	query := repo.db.Table("customers").Where(field+" = ?", value).Count(&count)
	if query.Error != nil {
		return false, query.Error
	}
	return count > 0, nil
}

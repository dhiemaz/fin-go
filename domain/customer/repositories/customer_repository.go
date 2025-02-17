package repositories

import (
	"context"
	"errors"
	"github.com/dhiemaz/fin-go/entities"
	"github.com/jinzhu/gorm"
)

// CustomerRepository interface
type CustomerRepository interface {
	Create(ctx context.Context, customer entities.Customer) error
	CreateBatch(ctx context.Context, customers []entities.Customer) error
	Update(ctx context.Context, customer entities.Customer) error
	Delete(ctx context.Context, customer entities.Customer) error
	GetAll(ctx context.Context, limit int, offset int) ([]entities.CustomerData, error)
	GetById(ctx context.Context, customerId int64) (entities.Customer, error)
	GetDataById(ctx context.Context, customerId int64) (entities.CustomerData, error)
	GetByUniqueId(ctx context.Context, uniqueId string) (entities.Customer, error)
	GetByDataUniqueId(ctx context.Context, uniqueId string) (entities.CustomerData, error)
	Count(ctx context.Context) (int64, error)
	ExistsRecord(ctx context.Context, field string, value string) (bool, error)
}

type Customer struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *Customer {
	return &Customer{
		db: db,
	}
}

// Create : create a customer
func (repo *Customer) Create(ctx context.Context, customer entities.Customer) error {
	result := repo.db.Create(&customer)
	return result.Error
}

// CreateBatch : create customer using batch mechanism
func (repo *Customer) CreateBatch(ctx context.Context, customers []entities.Customer) error {
	tx := repo.db.Begin()
	result := repo.db.Create(&customers)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	return tx.Commit().Error
}

// Update : update customer data
func (repo *Customer) Update(ctx context.Context, customer entities.Customer) error {
	result := repo.db.Save(&customer)
	return result.Error
}

// Delete : delete a customer
func (repo *Customer) Delete(ctx context.Context, customer entities.Customer) error {
	result := repo.db.Delete(&customer)
	return result.Error
}

// GetAll : get all customers
func (repo *Customer) GetAll(ctx context.Context, limit int, offset int) ([]entities.CustomerData, error) {
	var customers []entities.CustomerData
	result := repo.db.
		Table("view_customer_data").
		Limit(limit).Offset(offset).
		Find(&customers)
	return customers, result.Error
}

// GetById : get customer using id
func (repo *Customer) GetById(ctx context.Context, customerId int64) (entities.Customer, error) {
	var customer entities.Customer
	result := repo.db.Table("customers").First(&customer, customerId)
	if result.Error != nil {
		return entities.Customer{}, result.Error
	}
	return customer, result.Error
}

// GetDataById : get customer view data using id
func (repo *Customer) GetDataById(ctx context.Context, customerId int64) (entities.CustomerData, error) {
	var customer entities.CustomerData
	result := repo.db.Table("view_customer_data").First(&customer, customerId)
	if result.Error != nil {
		return entities.CustomerData{}, result.Error
	}
	return customer, result.Error
}

// GetByUniqueId : get customer data using unique id
func (repo *Customer) GetByUniqueId(ctx context.Context, uniqueId string) (entities.Customer, error) {
	var customer entities.Customer
	result := repo.db.Table("customers").First(&customer, "unique_id", uniqueId)
	if result.Error != nil {
		return entities.Customer{}, result.Error
	}
	return customer, result.Error
}

// GetByDataUniqueId : get customer view data using unique id
func (repo *Customer) GetByDataUniqueId(ctx context.Context, uniqueId string) (entities.CustomerData, error) {
	var customer entities.CustomerData
	result := repo.db.Table("view_customer_data").First(&customer, "unique_id", uniqueId)
	if result.Error != nil {
		return entities.CustomerData{}, result.Error
	}
	return customer, result.Error
}

// Count : get costumer data count
func (repo *Customer) Count(ctx context.Context) (int64, error) {
	var count int64
	result := repo.db.Table("customers").Count(&count)
	return count, result.Error
}

// ExistsRecord : check if record exist by valid fields
func (repo *Customer) ExistsRecord(ctx context.Context, field string, value string) (bool, error) {
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

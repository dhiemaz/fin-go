package usecase

import (
	"context"
	"fmt"
	"github.com/dhiemaz/fin-go/common/datetime"
	"github.com/dhiemaz/fin-go/common/encryption"
	"github.com/dhiemaz/fin-go/common/httputils"
	"github.com/dhiemaz/fin-go/domain/customer/repositories"
	"github.com/dhiemaz/fin-go/entities"
	"time"
)

// CustomerUseCase :
type CustomerUseCase interface {
	CreateCustomer(ctx context.Context, request entities.CreateCustomerRequest) error
	GetAllCustomers(ctx context.Context, params httputils.PaginationParams) ([]entities.CustomerData, int64, error)
	ChangeCustomerType(ctx context.Context, request entities.ChangeCustomerTypeRequest) error
	ChangeCustomerStatus(ctx context.Context, request entities.ChangeCustomerStatusRequest) error
	GetCustomerByUniqueId(ctx context.Context, uniqueId string) (entities.CustomerData, error)
	GetCustomerById(ctx context.Context, customerId int64) (entities.CustomerData, error)
	DeleteCustomer(ctx context.Context, customerId int64) error
	UpdateCustomerContacts(ctx context.Context, request entities.UpdateCustomerContactRequest) error
}

type Customer struct {
	Repository repositories.CustomerRepository
}

func NewCustomerUseCase(customerRepository repositories.CustomerRepository) *Customer {
	return &Customer{
		Repository: customerRepository,
	}
}

// CreateCustomer : create a new customer
func (customer *Customer) CreateCustomer(ctx context.Context, request entities.CreateCustomerRequest) error {
	if err := customer.checkDuplicatedValues(ctx, "identification_number", request.IdentificationNumber); err != nil {
		return httputils.NewConflictError(err.Error())
	}

	if err := customer.checkDuplicatedValues(ctx, "email", request.Email); err != nil {
		return httputils.NewConflictError(err.Error())
	}

	if err := customer.checkDuplicatedValues(ctx, "phone", request.Phone); err != nil {
		return httputils.NewConflictError(err.Error())
	}

	newCustomer := entities.Customer{
		CustomerType:         request.CustomerType,
		CustomerStatus:       entities.CustomerStatusPending,
		CustomerName:         request.CustomerName,
		IdentificationNumber: request.IdentificationNumber,
		Gender:               request.Gender,
		BirthDate:            datetime.StringToDate(request.BirthDate),
		Email:                request.Email,
		Phone:                request.Phone,
		Address:              request.Address,
		UniqueId:             encryption.GenerateUUID(),
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}
	return customer.Repository.Create(ctx, newCustomer)
}

// UpdateCustomerContacts : update customer contact data
func (customer *Customer) UpdateCustomerContacts(ctx context.Context, request entities.UpdateCustomerContactRequest) error {
	customerData, err := customer.Repository.GetById(ctx, request.CustomerId)
	if err != nil {
		return httputils.NewNotFoundError("Customer not found")
	}

	customerData.Email = request.Email
	customerData.Phone = request.Phone

	customer.Repository.Update(ctx, customerData)
	return nil
}

// GetAllCustomers : get all customers data
func (customer *Customer) GetAllCustomers(ctx context.Context, params httputils.PaginationParams) ([]entities.CustomerData, int64, error) {
	var customers []entities.CustomerData

	if params.CurrentPage < 0 || params.Limit < 1 {
		return nil, 0, httputils.NewBadRequestError("Incorrect current page or limit")
	}

	count, err := customer.Repository.Count(ctx)
	if err != nil {
		return customers, 0, httputils.NewNotFoundError("No customers found")
	}

	if count < 1 {
		return customers, count, httputils.NewNotFoundError("No customers found")
	}

	customers, err = customer.Repository.GetAll(ctx, params.Limit, params.CurrentPage)
	if err != nil {
		return customers, count, err
	}

	return customers, count, nil
}

// ChangeCustomerType : changing customer type
func (customer *Customer) ChangeCustomerType(ctx context.Context, request entities.ChangeCustomerTypeRequest) error {
	customerData, err := customer.Repository.GetById(ctx, request.CustomerId)
	if err != nil {
		return httputils.NewNotFoundError("Customer not found")
	}

	customerData.CustomerType = request.NewType
	customer.Repository.Update(ctx, customerData)
	return nil
}

// ChangeCustomerStatus : changing customer status
func (customer *Customer) ChangeCustomerStatus(ctx context.Context, request entities.ChangeCustomerStatusRequest) error {
	customerData, err := customer.Repository.GetById(ctx, request.CustomerId)
	if err != nil {
		return httputils.NewNotFoundError("Customer not found")
	}

	customerData.CustomerStatus = request.NewStatus
	customer.Repository.Update(ctx, customerData)
	return nil
}

// DeleteCustomer : delete a customer
func (customer *Customer) DeleteCustomer(ctx context.Context, customerId int64) error {
	customerData, err := customer.Repository.GetById(ctx, customerId)
	if err != nil {
		return httputils.NewNotFoundError("Customer not found")
	}

	customer.Repository.Delete(ctx, customerData)
	return nil
}

// GetCustomerById : get customer data using id
func (customer *Customer) GetCustomerById(ctx context.Context, customerId int64) (entities.CustomerData, error) {
	if customerId < 0 {
		return entities.CustomerData{}, httputils.NewBadRequestError("CustomerId must be greather than 0")
	}

	customerData, err := customer.Repository.GetDataById(ctx, customerId)
	if err != nil {
		return entities.CustomerData{}, httputils.NewNotFoundError("Customer not found")
	}
	return customerData, nil
}

// GetCustomerByUniqueId : get customer using unique id
func (customer *Customer) GetCustomerByUniqueId(ctx context.Context, uniqueId string) (entities.CustomerData, error) {
	if uniqueId == "" {
		return entities.CustomerData{}, httputils.NewBadRequestError("UniqueId cannot be null")
	}

	customerData, err := customer.Repository.GetByDataUniqueId(ctx, uniqueId)
	if err != nil {
		return entities.CustomerData{}, httputils.NewNotFoundError("Customer not found")
	}
	return customerData, nil
}

func (customer *Customer) checkDuplicatedValues(ctx context.Context, field string, value string) error {
	exists, err := customer.Repository.ExistsRecord(ctx, field, value)
	if err != nil {
		return fmt.Errorf("error checking %s existence: %w", field, err)
	}

	if exists {
		return fmt.Errorf("%s '%s' already exists", field, value)
	}
	return nil
}

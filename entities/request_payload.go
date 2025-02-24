package entities

// CreateCustomerRequest entity
type CreateCustomerRequest struct {
	CustomerType         CustomerType `json:"customer_type" validate:"required"`
	CustomerName         string       `json:"customer_name" validate:"required,min=2,max=150"`
	Gender               string       `json:"gender" validate:"required,oneof=Male Female Other"`
	BirthDate            string       `json:"birth_date" validate:"required"`
	IdentificationNumber string       `json:"identification_number" validate:"required,min=6,max=30"`
	Email                string       `json:"email,omitempty" validate:"omitempty,email,max=150"`
	Phone                string       `json:"phone,omitempty" validate:"omitempty,e164,max=20"`
	Address              string       `json:"address,omitempty" validate:"omitempty,max=200"`
}

// UpdateCustomerContactRequest entity
type UpdateCustomerContactRequest struct {
	CustomerId int64  `json:"customer_id" validate:"required"`
	Email      string `json:"email,omitempty" validate:"omitempty,email,max=150"`
	Phone      string `json:"phone,omitempty" validate:"omitempty,e164,max=20"`
}

// ChangeCustomerStatusRequest entity
type ChangeCustomerStatusRequest struct {
	CustomerId int64          `json:"customer_id" validate:"required"`
	NewStatus  CustomerStatus `json:"new_status" validate:"required"`
}

type ChangeCustomerTypeRequest struct {
	CustomerId int64        `json:"customer_id" validate:"required"`
	NewType    CustomerType `json:"new_type" validate:"required"`
}

type CreateAccountRequest struct {
	NickName   string  `json:"nick_name" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	CustomerID int64   `json:"customer_id" validate:"required"`
}

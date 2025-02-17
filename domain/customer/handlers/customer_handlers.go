package handlers

import (
	"bitbucket.org/rctiplus/almasbub"
	"fmt"
	"github.com/dhiemaz/fin-go/common/httputils"
	"github.com/dhiemaz/fin-go/common/serialization"
	"github.com/dhiemaz/fin-go/domain/customer/usecase"
	"github.com/dhiemaz/fin-go/entities"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	UseCase     usecase.CustomerUseCase
	infoLogger  *zap.Logger
	errorLogger *zap.Logger
}

func NewCustomerHandler(customerUseCase usecase.CustomerUseCase) *Handler {
	return &Handler{
		UseCase: customerUseCase,
		//infoLogger:  config.NewLogger("customers-info.log"),
		//errorLogger: config.NewLogger("customers-error.log"),
	}
}

func (customer *Handler) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := httputils.GetPaginationParams(r)
	customers, count, err := customer.UseCase.GetAllCustomers(ctx, params)
	if err != nil {
		httputils.HandleHTTPErrors(w, err)
		return
	}

	paginatedResult, err := httputils.NewPagination(r, customers, count, params.CurrentPage, params.Limit)
	if err != nil {
		httputils.HandleHTTPErrors(w, err)
		return
	}

	httputils.WriteJSONSimple(w, http.StatusOK, paginatedResult)
}

func (customer *Handler) getCustomerById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := almasbub.ToInt64(r.PathValue("id"))

	customerData, err := customer.UseCase.GetCustomerById(ctx, id)
	if err != nil {
		httputils.HandleHTTPErrors(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, customerData)
}

func (customer *Handler) getCustomerByUniqueId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uniqueId := r.PathValue("unique_id")

	customerData, err := customer.UseCase.GetCustomerByUniqueId(ctx, uniqueId)
	if err != nil {
		httputils.HandleHTTPErrors(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, customerData)
}

func (customer *Handler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var request entities.CreateCustomerRequest
	ctx := r.Context()

	if err := serialization.DecodeJson(r.Body, &request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	if err := httputils.Validate(request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	err := customer.UseCase.CreateCustomer(ctx, request)
	if err != nil {
		customer.errorLogger.Error(err.Error())
		httputils.HandleHTTPErrors(w, err)
		return
	}

	msg := fmt.Sprintf("Customer '%s' created", request.CustomerName)
	customer.infoLogger.Info(msg)
	httputils.WriteJSON(w, http.StatusCreated, msg)
}

func (customer *Handler) changeCustomerType(w http.ResponseWriter, r *http.Request) {
	var request entities.ChangeCustomerTypeRequest
	ctx := r.Context()

	if err := serialization.DecodeJson(r.Body, &request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	if err := httputils.Validate(request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	err := customer.UseCase.ChangeCustomerType(ctx, request)
	if err != nil {
		customer.errorLogger.Error(err.Error())
		httputils.HandleHTTPErrors(w, err)
		return
	}

	msg := fmt.Sprintf("Customer '%d' type changed", request.CustomerId)
	customer.infoLogger.Info(msg)
	httputils.WriteJSON(w, http.StatusCreated, msg)
}

func (customer *Handler) changeCustomerStatus(w http.ResponseWriter, r *http.Request) {
	var request entities.ChangeCustomerStatusRequest
	ctx := r.Context()

	if err := serialization.DecodeJson(r.Body, &request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	if err := httputils.Validate(request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	err := customer.UseCase.ChangeCustomerStatus(ctx, request)
	if err != nil {
		customer.errorLogger.Error(err.Error())
		httputils.HandleHTTPErrors(w, err)
		return
	}

	msg := fmt.Sprintf("Customer '%d' status changed", request.CustomerId)
	customer.infoLogger.Info(msg)
	httputils.WriteJSON(w, http.StatusCreated, msg)
}

func (customer *Handler) updateCustomerContacts(w http.ResponseWriter, r *http.Request) {
	var request entities.UpdateCustomerContactRequest
	ctx := r.Context()

	if err := serialization.DecodeJson(r.Body, &request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	if err := httputils.Validate(request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	err := customer.UseCase.UpdateCustomerContacts(ctx, request)
	if err != nil {
		httputils.HandleHTTPErrors(w, err)
		customer.errorLogger.Error(err.Error())
		return
	}

	msg := fmt.Sprintf("Customer '%d' contacts updated", request.CustomerId)
	customer.infoLogger.Info(msg)
	httputils.WriteJSON(w, http.StatusCreated, msg)
}

func (customer *Handler) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := almasbub.ToInt64(r.PathValue("id"))
	err := customer.UseCase.DeleteCustomer(r.Context(), customerId)
	if err != nil {
		customer.errorLogger.Error(err.Error())
		httputils.HandleHTTPErrors(w, err)
		return
	}

	msg := fmt.Sprintf("Customer '%d' deleted", customerId)
	customer.infoLogger.Info(msg)
	httputils.WriteJSON(w, http.StatusNoContent, msg)
}

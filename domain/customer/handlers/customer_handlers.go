package handlers

import (
	"bitbucket.org/rctiplus/almasbub"
	"fmt"
	"github.com/dhiemaz/fin-go/common"
	"github.com/dhiemaz/fin-go/entities"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"net/http"
)

type CustomerHandler struct {
	//service *services.CustomerService
	infoLogger  *zap.Logger
	errorLogger *zap.Logger
}

func NewCustomerHandler(db *gorm.DB) *CustomerHandler {
	return &CustomerHandler{
		//service:     services.NewCustomerService(db),
		//infoLogger:  config.NewLogger("customers-info.log"),
		//errorLogger: config.NewLogger("customers-error.log"),
	}
}

func (handler *CustomerHandler) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	params := common.GetPaginationParams(r)
	customers, err := handler.service.GetAllCustomers(r, r.Context(), params)
	if err != nil {
		common.HandleHTTPErrors(w, err)
		return
	}
	common.WriteJSONSimple(w, http.StatusOK, customers)
}

func (handler *CustomerHandler) getCustomerById(w http.ResponseWriter, r *http.Request) {
	id := almasbub.ToInt64(r.PathValue("id"))
	customer, err := handler.service.GetCustomerById(r.Context(), id)
	if err != nil {
		common.HandleHTTPErrors(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, customer)
}

func (handler *CustomerHandler) getCustomerByUniqueId(w http.ResponseWriter, r *http.Request) {
	uniqueId := r.PathValue("unique_id")
	customer, err := handler.service.GetCustomerByUniqueId(r.Context(), uniqueId)
	if err != nil {
		common.HandleHTTPErrors(w, err)
		return
	}
	common.WriteJSON(w, http.StatusOK, customer)
}

func (handler *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var request entities.CreateCustomerRequest
	err := handler.service.CreateCustomer(r, r.Context(), request)
	if err != nil {
		common.HandleHTTPErrors(w, err)
		handler.errorLogger.Error(err.Error())
		return
	}

	msg := fmt.Sprintf("Customer '%s' created", request.CustomerName)
	handler.infoLogger.Info(msg)
	common.WriteJSON(w, http.StatusCreated, msg)
}

func (handler *CustomerHandler) changeCustomerType(w http.ResponseWriter, r *http.Request) {
	var request entities.ChangeCustomerTypeRequest
	err := handler.service.ChangeCustomerType(r, r.Context(), request)
	if err != nil {
		common.HandleHTTPErrors(w, err)
		handler.errorLogger.Error(err.Error())
		return
	}
	msg := fmt.Sprintf("Customer '%d' type changed", request.CustomerId)
	handler.infoLogger.Info(msg)
	common.WriteJSON(w, http.StatusCreated, msg)
}

func (handler *CustomerHandler) changeCustomerStatus(w http.ResponseWriter, r *http.Request) {
	var request entities.ChangeCustomerStatusRequest
	err := handler.service.ChangeCustomerStatus(r, r.Context(), request)
	if err != nil {
		common.HandleHTTPErrors(w, err)
		handler.errorLogger.Error(err.Error())
		return
	}
	msg := fmt.Sprintf("Customer '%d' status changed", request.CustomerId)
	handler.infoLogger.Info(msg)
	common.WriteJSON(w, http.StatusCreated, msg)
}

func (handler *CustomerHandler) updateCustomerContacts(w http.ResponseWriter, r *http.Request) {
	var request entities.UpdateCustomerContactRequest
	err := handler.service.UpdateCustomerContacts(r, r.Context(), request)
	if err != nil {
		common.HandleHTTPErrors(w, err)
		handler.errorLogger.Error(err.Error())
		return
	}
	msg := fmt.Sprintf("Customer '%d' contacts updated", request.CustomerId)
	handler.infoLogger.Info(msg)
	common.WriteJSON(w, http.StatusCreated, msg)
}

func (handler *CustomerHandler) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := almasbub.ToInt64(r.PathValue("id"))
	err := handler.service.DeleteCustomer(r.Context(), customerId)
	if err != nil {
		common.HandleHTTPErrors(w, err)
		handler.errorLogger.Error(err.Error())
		return
	}
	msg := fmt.Sprintf("Customer '%d' deleted", customerId)
	handler.infoLogger.Info(msg)
	common.WriteJSON(w, http.StatusNoContent, msg)
}

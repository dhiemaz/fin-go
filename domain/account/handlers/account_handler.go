package handlers

import (
	"bitbucket.org/rctiplus/almasbub"
	"fmt"
	"github.com/dhiemaz/fin-go/common/httputils"
	"github.com/dhiemaz/fin-go/common/serialization"
	"github.com/dhiemaz/fin-go/domain/account/usecase"
	"github.com/dhiemaz/fin-go/entities"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	UseCase     usecase.AccountUseCase
	infoLogger  *zap.Logger
	errorLogger *zap.Logger
}

func NewAccountHandler(accountUseCase usecase.AccountUseCase) *Handler {
	return &Handler{
		UseCase: accountUseCase,
		//infoLogger:  config.NewLogger("customers-info.log"),
		//errorLogger: config.NewLogger("customers-error.log"),
	}
}

func (account *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	var request entities.CreateAccountRequest
	ctx := r.Context()

	if err := serialization.DecodeJson(r.Body, &request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	if err := httputils.Validate(request); err != nil {
		httputils.HandleHTTPErrors(w, httputils.NewBadRequestError(err.Error()))
	}

	err := account.UseCase.CreateAccount(ctx, request)
	if err != nil {
		account.errorLogger.Error(err.Error())
		httputils.HandleHTTPErrors(w, err)
		return
	}

	msg := fmt.Sprintf("Account '%s' created", request.NickName)
	account.infoLogger.Info(msg)
	httputils.WriteJSON(w, http.StatusCreated, msg)
}

func (account *Handler) deleteAccount(w http.ResponseWriter, r *http.Request) {
	accountId := almasbub.ToInt64(r.PathValue("id"))
	err := account.UseCase.DeleteAccount(r.Context(), accountId)
	if err != nil {
		account.errorLogger.Error(err.Error())
		httputils.HandleHTTPErrors(w, err)
		return
	}

	msg := fmt.Sprintf("Account '%d' deleted", accountId)
	account.infoLogger.Info(msg)
	httputils.WriteJSON(w, http.StatusNoContent, msg)
}

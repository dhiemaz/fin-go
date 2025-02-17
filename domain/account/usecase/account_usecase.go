package usecase

import (
	"context"
	"fmt"
	"github.com/dhiemaz/fin-go/common/encryption"
	"github.com/dhiemaz/fin-go/common/httputils"
	"github.com/dhiemaz/fin-go/domain/account/repositories"
	"github.com/dhiemaz/fin-go/entities"
	"time"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, request entities.CreateAccountRequest) error
	GetAllAccounts(ctx context.Context, params httputils.PaginationParams) ([]entities.Account, int64, error)
	GetAccountByCIF(ctx context.Context, uniqueId string) (entities.Account, error)
	GetAccountById(ctx context.Context, accountId int64) (entities.Account, error)
	DeleteAccount(ctx context.Context, accountId int64) error
}

type Account struct {
	Repository repositories.AccountRepository
}

func NewAccountUseCase(accountRepository repositories.AccountRepository) *Account {
	return &Account{
		Repository: accountRepository,
	}
}

func (account *Account) CreateAccount(ctx context.Context, request entities.CreateAccountRequest) error {
	cif := encryption.GenerateCIF()

	if err := account.checkDuplicatedValues(ctx, "cif", cif); err != nil {
		return httputils.NewConflictError(err.Error())
	}

	newAccount := entities.Account{
		CIF:        cif,
		NickName:   request.NickName,
		Amount:     request.Amount,
		CustomerID: request.CustomerID,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	return account.Repository.Create(ctx, newAccount)
}

func (account *Account) GetAllAccounts(ctx context.Context, params httputils.PaginationParams) ([]entities.Account, int64, error) {
	var accounts []entities.Account

	if params.CurrentPage < 0 || params.Limit < 1 {
		return nil, 0, httputils.NewBadRequestError("Incorrect current page or limit")
	}

	count, err := account.Repository.Count(ctx)
	if err != nil {
		return accounts, 0, httputils.NewNotFoundError("No customers found")
	}

	if count < 1 {
		return accounts, count, httputils.NewNotFoundError("No customers found")
	}

	accounts, err = account.Repository.GetAll(ctx, params.Limit, params.CurrentPage)
	if err != nil {
		return accounts, count, err
	}

	return accounts, count, nil
}

func (account *Account) GetAccountByCIF(ctx context.Context, cif string) (entities.Account, error) {
	if cif == "" {
		return entities.Account{}, httputils.NewBadRequestError("CIF cannot be null")
	}

	accountData, err := account.Repository.GetByCIF(ctx, cif)
	if err != nil {
		return entities.Account{}, httputils.NewNotFoundError("Account not found")
	}
	return accountData, nil
}

func (account *Account) DeleteAccount(ctx context.Context, accountId int64) error {
	accountData, err := account.Repository.GetDataById(ctx, accountId)
	if err != nil {
		return httputils.NewNotFoundError("Account not found")
	}

	account.Repository.Delete(ctx, accountData)
	return nil
}

func (account *Account) checkDuplicatedValues(ctx context.Context, field string, value string) error {
	exists, err := account.Repository.ExistsRecord(ctx, field, value)
	if err != nil {
		return fmt.Errorf("error checking %s existence: %w", field, err)
	}

	if exists {
		return fmt.Errorf("%s '%s' already exists", field, value)
	}
	return nil
}

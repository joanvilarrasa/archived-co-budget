package app

import (
	"co-budget/data"
	"co-budget/lib"
)

type Account struct {
	ID             int64
	Name           string
	Description    string
	CreatedAt      string
	InitialBalance string
	Type           string
}

type AccountsProps struct {
	Accounts []data.Account
	Error    string
}

func Accounts(errorMessage string) string {
	props := AccountsProps{Error: errorMessage, Accounts: make([]data.Account, 0)}

	if data.AccountsStore == nil {
		props.Error = "Accounts store not initialized"
		return lib.ParseHtmlTemplate("./app/accounts.html", props)
	}

	accounts, err := data.AccountsStore.List()
	if err != nil {
		if props.Error == "" {
			props.Error = "Failed to load accounts"
		}
		return lib.ParseHtmlTemplate("./app/accounts.html", props)
	}

	props.Accounts = make([]data.Account, 0, len(accounts))
	for _, account := range accounts {
		props.Accounts = append(props.Accounts, data.Account{
			ID:             account.ID,
			Name:           account.Name,
			Description:    account.Description,
			CreatedAt:      account.CreatedAt,
			InitialBalance: account.InitialBalance,
			Type:           account.Type,
		})
	}

	return lib.ParseHtmlTemplate("./app/accounts.html", props)
}

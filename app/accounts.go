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

func Accounts() string {
	return AccountsWithError("")
}

func AccountsWithError(requestError string) string {

	accounts, accountsRes := data.AccountGetAll()
	errorMessage := requestError
	if accountsRes != data.AS_Ok && errorMessage == "" {
		errorMessage = "Error retrieving accounts"
	}

	props := AccountsProps{Error: errorMessage, Accounts: accounts}

	return lib.ParseHtmlTemplate("./app/accounts.html", props)
}

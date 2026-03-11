package app

import "co-budget/lib"

func Accounts() string {
	return lib.ParseHtmlTemplate("./app/accounts.html", nil)
}

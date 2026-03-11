package app

import "co-budget/lib"

func Budgets() string {
	return lib.ParseHtmlTemplate("./app/budgets.html", nil)
}

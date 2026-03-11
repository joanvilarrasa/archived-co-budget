package app

import "co-budget/lib"

func Expenses() string {
	return lib.ParseHtmlTemplate("./app/expenses.html", nil)
}

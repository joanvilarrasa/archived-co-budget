package app

import "co-budget/lib"

func navHeader() string {
	return "CO-Budget"
}

type HomeProps struct {
	Dashboard   string
	Accounts    string
	Budgets     string
	Expenses    string
	InitialPage string
}

func Layout() string {
	layoutdata := HomeProps{
		Dashboard:   Dashboard(),
		Accounts:    Accounts(),
		Budgets:     Budgets(),
		Expenses:    Expenses(),
		InitialPage: "dashboard",
	}
	return lib.ParseHtmlTemplate("./app/layout.html", layoutdata)
}

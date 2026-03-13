package app

import "co-budget/lib"

func navHeader() string {
	return "CO-Budget"
}

type HomeProps struct {
	DashboardPage string
	AccountsPage  string
	BudgetsPage   string
	ExpensesPage  string
}

func Layout() string {
	layoutdata := HomeProps{
		DashboardPage: Dashboard(),
		AccountsPage:  Accounts(),
		BudgetsPage:   Budgets(),
		ExpensesPage:  Expenses(),
	}
	return lib.ParseHtmlTemplate("./app/layout.html", layoutdata)
}

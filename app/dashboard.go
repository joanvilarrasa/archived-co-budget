package app

import "co-budget/lib"

func Dashboard() string {
	return lib.ParseHtmlTemplate("./app/dashboard.html", nil)
}

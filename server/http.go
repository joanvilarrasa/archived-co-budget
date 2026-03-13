package server

import (
	"co-budget/app"
	"co-budget/data"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HTTPServer struct {
}

func NewHTTPServer() *http.Server {
	srv := &HTTPServer{}
	mux := http.NewServeMux()

	mux.HandleFunc("/", srv.firstRender)
	mux.HandleFunc("/accounts", srv.createAccount)
	mux.HandleFunc("/accounts/delete", srv.deleteAccount)
	mux.HandleFunc("/datastar.js", datastarJS)
	mux.HandleFunc("/main.css", mainCss)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func (s *HTTPServer) firstRender(w http.ResponseWriter, r *http.Request) {
	initialPage := sanitizePage(r.URL.Query().Get("page"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, app.Layout(r.URL.Query().Get("error"), initialPage))
}

func (s *HTTPServer) createAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if parseErr := r.ParseForm(); parseErr != nil {
		redirectAccounts(w, r, "Invalid form body")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))
	accountType := strings.TrimSpace(r.FormValue("type"))
	balanceRaw := strings.TrimSpace(r.FormValue("initial_balance"))

	if name == "" || description == "" || balanceRaw == "" {
		redirectAccounts(w, r, "Name, description and initial balance are required")
		return
	}

	if !isValidAccountType(accountType) {
		redirectAccounts(w, r, "Type must be one of LTB, MTB or STB")
		return
	}

	initialBalance, parseErr := strconv.ParseFloat(balanceRaw, 64)
	if parseErr != nil {
		redirectAccounts(w, r, "Initial balance must be a number")
		return
	}

	if data.AccountsStore == nil {
		redirectAccounts(w, r, "Accounts store not initialized")
		return
	}

	if createErr := data.AccountsStore.Create(name, description, initialBalance, accountType); createErr != nil {
		redirectAccounts(w, r, "Failed to create account")
		return
	}

	redirectAccounts(w, r, "")
}

func (s *HTTPServer) deleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if parseErr := r.ParseForm(); parseErr != nil {
		redirectAccounts(w, r, "Invalid delete request")
		return
	}

	idRaw := strings.TrimSpace(r.FormValue("id"))
	id, parseErr := strconv.ParseInt(idRaw, 10, 64)
	if parseErr != nil {
		redirectAccounts(w, r, "Invalid account id")
		return
	}

	if data.AccountsStore == nil {
		redirectAccounts(w, r, "Accounts store not initialized")
		return
	}

	if deleteErr := data.AccountsStore.Delete(id); deleteErr != nil {
		redirectAccounts(w, r, "Failed to delete account")
		return
	}

	redirectAccounts(w, r, "")
}

func datastarJS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "datastar.js")
}

func mainCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	http.ServeFile(w, r, "main.css")
}

func sanitizePage(page string) string {
	switch page {
	case "accounts", "dashboard", "expenses", "budgets":
		return page
	default:
		return "dashboard"
	}
}

func redirectAccounts(w http.ResponseWriter, r *http.Request, message string) {
	redirectURL := "/?page=accounts"
	if message != "" {
		redirectURL = redirectURL + "&error=" + url.QueryEscape(message)
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func isValidAccountType(accountType string) bool {
	return accountType == "LTB" || accountType == "MTB" || accountType == "STB"
}

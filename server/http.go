package server

import (
	"co-budget/app"
	"co-budget/data"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/starfederation/datastar-go/datastar"
)

type HttpServerResponseDTO struct {
	Succes  bool   `json:"succes"`
	Message string `json:"message"`
}

type HTTPServer struct {
}

func NewHTTPServer() *http.Server {
	srv := &HTTPServer{}
	mux := http.NewServeMux()

	mux.HandleFunc("/", srv.firstRender)
	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		statusCode, response := srv.createAccount(w, r)
		writeJsonResponse(w, statusCode, response)
	})
	mux.HandleFunc("/accounts/delete", srv.deleteAccount)
	mux.HandleFunc("/datastar.js", datastarJS)
	mux.HandleFunc("/main.css", mainCss)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, response HttpServerResponseDTO) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func (s *HTTPServer) firstRender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, app.Layout())
}

func (s *HTTPServer) createAccount(w http.ResponseWriter, r *http.Request) (int, HttpServerResponseDTO) {
	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, HttpServerResponseDTO{
			Succes:  false,
			Message: "method_not_allowed",
		}
	}

	if parseErr := r.ParseForm(); parseErr != nil {
		return http.StatusBadRequest, HttpServerResponseDTO{
			Succes:  false,
			Message: "The form body is invalid",
		}
	}

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))
	accountType := strings.TrimSpace(r.FormValue("type"))
	balanceRaw := strings.TrimSpace(r.FormValue("initial_balance"))

	if name == "" || description == "" || balanceRaw == "" {
		return http.StatusBadRequest, HttpServerResponseDTO{
			Succes:  false,
			Message: "Name, description and initial balance are required fields",
		}
	}

	if !isValidAccountType(accountType) {
		return http.StatusBadRequest, HttpServerResponseDTO{
			Succes:  false,
			Message: "The account type is invalid",
		}
	}

	initialBalance, parseErr := strconv.ParseFloat(balanceRaw, 64)
	if parseErr != nil {
		return http.StatusBadRequest, HttpServerResponseDTO{
			Succes:  false,
			Message: "Initial balance must be a number",
		}
	}

	createRes := data.AccountCreate(name, description, initialBalance, accountType)
	switch createRes {
	case data.AS_Ok:
		s.patchAccounts(w, r)
		return http.StatusAccepted, HttpServerResponseDTO{
			Succes:  true,
			Message: "Created successfully",
		}

	default:
		return http.StatusInternalServerError, HttpServerResponseDTO{
			Succes:  false,
			Message: "Some uncontrolled error has ocurred",
		}
	}

}

func (s *HTTPServer) deleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if parseErr := r.ParseForm(); parseErr != nil {
		s.patchAccounts(w, r)
		return
	}

	idRaw := strings.TrimSpace(r.FormValue("id"))
	id, parseErr := strconv.ParseInt(idRaw, 10, 64)
	if parseErr != nil {
		s.patchAccounts(w, r)
		return
	}

	if deleteRes := data.AccountDelete(id); deleteRes != data.AS_Ok {
		s.patchAccounts(w, r)
		return
	}

	s.patchAccounts(w, r)
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

func isValidAccountType(accountType string) bool {
	return accountType == "LTB" || accountType == "MTB" || accountType == "STB"
}

func (s *HTTPServer) patchAccounts(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(app.Accounts())
}

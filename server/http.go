package server

import (
	"co-budget/app"
	"co-budget/data"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

type HttpServerResponseDTO struct {
	Succes  bool   `json:"succes"`
	Message string `json:"message"`
}

type HTTPServer struct {
	mu             sync.Mutex
	sseConnections map[int64]sseConnection
	nextSSEID      int64
}

type sseConnection struct {
	stream interface {
		PatchElements(string, ...datastar.PatchElementOption) error
	}
	done <-chan struct{}
}

func NewHTTPServer() *http.Server {
	srv := &HTTPServer{
		sseConnections: map[int64]sseConnection{},
	}
	mux := http.NewServeMux()

	go srv.broadcastSSEStatus()

	mux.HandleFunc("/", srv.firstRender)
	mux.HandleFunc("/sse", srv.sse)
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

func (s *HTTPServer) sse(w http.ResponseWriter, r *http.Request) {
	stream := datastar.NewSSE(w, r)
	connectionID := s.addSSEConnection(stream, r.Context().Done())
	defer s.removeSSEConnection(connectionID)

	<-r.Context().Done()
}

func (s *HTTPServer) addSSEConnection(stream interface {
	PatchElements(string, ...datastar.PatchElementOption) error
}, done <-chan struct{}) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextSSEID++
	connectionID := s.nextSSEID
	s.sseConnections[connectionID] = sseConnection{stream: stream, done: done}

	return connectionID
}

func (s *HTTPServer) removeSSEConnection(connectionID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sseConnections, connectionID)
}

func (s *HTTPServer) broadcastSSEStatus() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		s.broadcastToSSEConnections(fmt.Sprintf(`<div id="sse-status">SSE Status ON %s</div>`, t.Format("15:04:05")))
	}
}

func (s *HTTPServer) broadcastToSSEConnections(message string) {
	s.mu.Lock()
	activeConnections := make([]sseConnection, 0, len(s.sseConnections))

	for connectionID, connection := range s.sseConnections {
		select {
		case <-connection.done:
			delete(s.sseConnections, connectionID)
		default:
			activeConnections = append(activeConnections, connection)
		}
	}

	s.mu.Unlock()

	for _, connection := range activeConnections {
		_ = connection.stream.PatchElements(message)
	}
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
		s.broadcastToSSEConnections(app.Accounts())
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

	s.broadcastToSSEConnections(app.Accounts())
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

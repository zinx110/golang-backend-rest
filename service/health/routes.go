package health

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zinx110/golang-backend-rest/utils"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/healthz", h.handleHealthz)
}

func (h *Handler) handleHealthz(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}

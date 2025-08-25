package health

import (
	"database/sql"
	"fmt"
	"log"
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
	err := h.db.Ping()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error connecting the db %w", err))
		log.Fatal(err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}

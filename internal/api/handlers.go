package api

import (
	"encoding/json"
	"example.com/pz3-http/internal/storage"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Handlers struct{ Store *storage.MemoryStore }

func NewHandlers(store *storage.MemoryStore) *Handlers { return &Handlers{Store: store} }
func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		filtered := tasks[:0]
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), strings.ToLower(q)) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}
	JSON(w, http.StatusOK, tasks)
}
func isJSON(w http.ResponseWriter, r *http.Request) bool {
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		BadRequest(w, "Content-Type must be application/json")
		return false
	}
	return true
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	if !isJSON(w, r) {
		return
	}
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json")
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	length := utf8.RuneCountInString(req.Title)
	if length < 3 || length > 140 {
		JSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: "title must contain 3 to 140 characters"})
		return
	}
	JSON(w, http.StatusCreated, h.Store.Create(req.Title))
}
func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		BadRequest(w, "invalid id")
		return 0, false
	}
	return id, true
}
func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	t, err := h.Store.Get(id)
	if err != nil {
		NotFound(w, "task not found")
		return
	}
	JSON(w, http.StatusOK, t)
}

type patchTaskRequest struct {
	Done *bool `json:"done"`
}

func (h *Handlers) PatchTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if !isJSON(w, r) {
		return
	}
	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Done == nil {
		BadRequest(w, "require boolean done")
		return
	}
	t, err := h.Store.SetDone(id, *req.Done)
	if err != nil {
		NotFound(w, "task not found")
		return
	}
	JSON(w, http.StatusOK, t)
}
func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := h.Store.Delete(id); err != nil {
		NotFound(w, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

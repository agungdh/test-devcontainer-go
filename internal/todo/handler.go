package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"todoapp/pkg/response"
)

type Handler struct {
	svc      Service
	validate *validator.Validate
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc, validate: validator.New()}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "body JSON tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	t, err := h.svc.Create(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, t)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var isCompleted *bool
	if v := q.Get("is_completed"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "bad_request", "is_completed harus true/false")
			return
		}
		isCompleted = &b
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	todos, err := h.svc.List(r.Context(), ListFilter{
		Q:           q.Get("q"),
		IsCompleted: isCompleted,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	if todos == nil {
		todos = []Todo{}
	}
	response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "id harus angka")
		return
	}
	t, err := h.svc.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		response.Error(w, http.StatusNotFound, "not_found", "todo tidak ditemukan")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, t)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "id harus angka")
		return
	}
	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "body JSON tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	t, err := h.svc.Update(r.Context(), id, req)
	if errors.Is(err, ErrNotFound) {
		response.Error(w, http.StatusNotFound, "not_found", "todo tidak ditemukan")
		return
	}
	if errors.Is(err, ErrEmptyUpdate) {
		response.Error(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, t)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "bad_request", "id harus angka")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(w, http.StatusNotFound, "not_found", "todo tidak ditemukan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

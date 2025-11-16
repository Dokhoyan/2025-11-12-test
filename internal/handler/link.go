package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type AddLinksRequest struct {
	URLs []string `json:"urls"`
}

type AddLinksResponse struct {
	ID    int64          `json:"id"`
	Links []LinkResponse `json:"links"`
}

type LinkResponse struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

func (h *Handler) AddLinks(w http.ResponseWriter, r *http.Request) {
	var req AddLinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrorResponse("Invalid request body", w, http.StatusBadRequest)
		return
	}

	if err := h.validateAddLinksRequest(&req); err != nil {
		newErrorResponse(err.Error(), w, http.StatusBadRequest)
		return
	}

	linkSet, err := h.service.AddLinks(r.Context(), req.URLs)
	if err != nil {
		newErrorResponse(err.Error(), w, http.StatusInternalServerError)
		return
	}

	response := AddLinksResponse{
		ID:    linkSet.ID,
		Links: make([]LinkResponse, len(linkSet.Links)),
	}

	for i, link := range linkSet.Links {
		response.Links[i] = LinkResponse{
			URL:    link.URL,
			Status: link.Status,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		newErrorResponse("Failed to encode response", w, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetLinkSet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		newErrorResponse("id parameter is required", w, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		newErrorResponse("Invalid id parameter", w, http.StatusBadRequest)
		return
	}

	linkSet, err := h.service.GetLinkSet(id)
	if err != nil {
		newErrorResponse(err.Error(), w, http.StatusNotFound)
		return
	}

	response := AddLinksResponse{
		ID:    linkSet.ID,
		Links: make([]LinkResponse, len(linkSet.Links)),
	}

	for i, link := range linkSet.Links {
		response.Links[i] = LinkResponse{
			URL:    link.URL,
			Status: link.Status,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		newErrorResponse("Failed to encode response", w, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) validateAddLinksRequest(req *AddLinksRequest) error {
	if len(req.URLs) == 0 {
		return &ValidationError{Message: "URLs list is empty"}
	}

	for _, url := range req.URLs {
		if url == "" {
			return &ValidationError{Message: "URL cannot be empty"}
		}
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

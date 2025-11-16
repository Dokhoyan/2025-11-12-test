package handler

import (
	"encoding/json"
	"net/http"
)

type GenerateReportRequest struct {
	LinksNum []int64 `json:"links_num"`
}

func (h *Handler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrorResponse("Invalid request body", w, http.StatusBadRequest)
		return
	}

	if err := h.validateGenerateReportRequest(&req); err != nil {
		newErrorResponse(err.Error(), w, http.StatusBadRequest)
		return
	}

	pdfData, err := h.service.GenerateReport(req.LinksNum)
	if err != nil {
		newErrorResponse(err.Error(), w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=report.pdf")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(pdfData); err != nil {
		newErrorResponse("Failed to write PDF data", w, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) validateGenerateReportRequest(req *GenerateReportRequest) error {
	if len(req.LinksNum) == 0 {
		return &ValidationError{Message: "links_num list is empty"}
	}

	for _, id := range req.LinksNum {
		if id <= 0 {
			return &ValidationError{Message: "link number must be positive"}
		}
	}

	return nil
}

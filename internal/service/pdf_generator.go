package service

import (
	"bytes"
	"fmt"

	"github.com/Dokhoyan/2025-11-12-test/internal/domain"

	"codeberg.org/go-pdf/fpdf"
)

// PDFGenerator генерирует PDF отчеты
type PDFGenerator interface {
	GenerateReport(linkSets []*domain.LinkSet) ([]byte, error)
}

// FPDFGenerator реализует генерацию PDF с помощью fpdf
type FPDFGenerator struct{}

// NewFPDFGenerator создает новый PDF генератор
func NewFPDFGenerator() *FPDFGenerator {
	return &FPDFGenerator{}
}

// GenerateReport генерирует PDF отчет по наборам ссылок
func (g *FPDFGenerator) GenerateReport(linkSets []*domain.LinkSet) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Link Status Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)

	for _, linkSet := range linkSets {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Link Set #%d", linkSet.ID))
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)
		for _, link := range linkSet.Links {
			statusText := link.Status
			if statusText == "available" {
				pdf.SetTextColor(0, 128, 0)
			} else {
				pdf.SetTextColor(255, 0, 0)
			}
			pdf.Cell(40, 8, fmt.Sprintf("  %s - %s", link.URL, statusText))
			pdf.Ln(6)
			pdf.SetTextColor(0, 0, 0)
		}
		pdf.Ln(8)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

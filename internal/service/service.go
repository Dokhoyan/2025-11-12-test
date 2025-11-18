package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Dokhoyan/2025-11-12-test/internal/domain"
	"github.com/Dokhoyan/2025-11-12-test/internal/repository"
)

type Service struct {
	repo         repository.Repository
	checker      LinkChecker
	pdfGenerator PDFGenerator
	timeout      time.Duration
}

func NewService(repo repository.Repository, checker LinkChecker, pdfGenerator PDFGenerator, timeout time.Duration) *Service {
	return &Service{
		repo:         repo,
		checker:      checker,
		pdfGenerator: pdfGenerator,
		timeout:      timeout,
	}
}

func (s *Service) AddLinks(ctx context.Context, urls []string) (*domain.LinkSet, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("urls list is empty")
	}

	links, err := s.checker.CheckLinks(ctx, urls)
	if err != nil {
		return nil, fmt.Errorf("failed to check links: %w", err)
	}

	linkSet := &domain.LinkSet{
		Links: links,
	}

	if err := s.repo.SaveLinkSet(linkSet); err != nil {
		return nil, fmt.Errorf("failed to save link set: %w", err)
	}

	return linkSet, nil
}

func (s *Service) GetLinkSet(id int64) (*domain.LinkSet, error) {
	return s.repo.GetLinkSet(id)
}

func (s *Service) GenerateReport(linkSetIDs []int64) ([]byte, error) {
	if len(linkSetIDs) == 0 {
		return nil, fmt.Errorf("link set IDs list is empty")
	}

	linkSets, err := s.repo.GetLinkSets(linkSetIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get link sets: %w", err)
	}

	if len(linkSets) == 0 {
		return nil, fmt.Errorf("no link sets found")
	}

	pdfData, err := s.pdfGenerator.GenerateReport(linkSets)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfData, nil
}

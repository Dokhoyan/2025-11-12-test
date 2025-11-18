package repository

import "github.com/Dokhoyan/2025-11-12-test/internal/domain"

type Repository interface {
	SaveLinkSet(linkSet *domain.LinkSet) error
	GetLinkSet(id int64) (*domain.LinkSet, error)
	GetLinkSets(ids []int64) ([]*domain.LinkSet, error)
	GetAllLinkSets() ([]*domain.LinkSet, error)
	Close() error
}

package repository

import "github.com/Dokhoyan/2025-11-12-test/internal/domain"

// Repository определяет интерфейс для работы с хранилищем данных
type Repository interface {
	// SaveLinkSet сохраняет набор ссылок и возвращает присвоенный номер
	SaveLinkSet(linkSet *domain.LinkSet) error

	// GetLinkSet получает набор ссылок по номеру
	GetLinkSet(id int64) (*domain.LinkSet, error)

	// GetLinkSets получает наборы ссылок по списку номеров
	GetLinkSets(ids []int64) ([]*domain.LinkSet, error)

	// GetAllLinkSets получает все наборы ссылок
	GetAllLinkSets() ([]*domain.LinkSet, error)

	// Close закрывает репозиторий и сохраняет данные
	Close() error
}

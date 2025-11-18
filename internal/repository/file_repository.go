package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Dokhoyan/2025-11-12-test/internal/domain"
)

type FileRepository struct {
	mu        sync.RWMutex
	linkSets  map[int64]*domain.LinkSet
	nextID    int64
	filePath  string
	saveMutex sync.Mutex
}

func NewFileRepository(filePath string) (*FileRepository, error) {
	repo := &FileRepository{
		linkSets: make(map[int64]*domain.LinkSet),
		nextID:   1,
		filePath: filePath,
	}

	if err := repo.load(); err != nil {
		return nil, fmt.Errorf("failed to load repository: %w", err)
	}

	return repo, nil
}

func (r *FileRepository) SaveLinkSet(linkSet *domain.LinkSet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if linkSet.ID == 0 {
		linkSet.ID = r.nextID
		r.nextID++
	}
	linkSet.CreatedAt = time.Now()

	r.linkSets[linkSet.ID] = linkSet

	return r.save()
}

func (r *FileRepository) GetLinkSet(id int64) (*domain.LinkSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	linkSet, exists := r.linkSets[id]
	if !exists {
		return nil, fmt.Errorf("link set with id %d not found", id)
	}

	return linkSet, nil
}

func (r *FileRepository) GetLinkSets(ids []int64) ([]*domain.LinkSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.LinkSet
	for _, id := range ids {
		if linkSet, exists := r.linkSets[id]; exists {
			result = append(result, linkSet)
		}
	}

	return result, nil
}

func (r *FileRepository) GetAllLinkSets() ([]*domain.LinkSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.LinkSet, 0, len(r.linkSets))
	for _, linkSet := range r.linkSets {
		result = append(result, linkSet)
	}

	return result, nil
}

func (r *FileRepository) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.save()
}

func (r *FileRepository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var storage struct {
		LinkSets map[int64]*domain.LinkSet `json:"link_sets"`
		NextID   int64                     `json:"next_id"`
	}

	if err := json.Unmarshal(data, &storage); err != nil {
		return err
	}

	r.linkSets = storage.LinkSets
	if r.linkSets == nil {
		r.linkSets = make(map[int64]*domain.LinkSet)
	}
	r.nextID = storage.NextID
	if r.nextID == 0 {
		r.nextID = 1
	}

	return nil
}

func (r *FileRepository) save() error {
	r.saveMutex.Lock()
	defer r.saveMutex.Unlock()

	storage := struct {
		LinkSets map[int64]*domain.LinkSet `json:"link_sets"`
		NextID   int64                     `json:"next_id"`
	}{
		LinkSets: r.linkSets,
		NextID:   r.nextID,
	}

	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := r.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0600); err != nil {
		return err
	}

	return os.Rename(tmpFile, r.filePath)
}

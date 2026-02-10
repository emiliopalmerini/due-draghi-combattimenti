package monster

import (
	"github.com/emiliopalmerini/due-draghi-combattimenti/internal/domain/monster"
)

// Service provides monster search use cases.
type Service struct {
	repo monster.Repository
}

// NewService creates a new monster application service.
func NewService(repo monster.Repository) *Service {
	return &Service{repo: repo}
}

// SearchMonsters returns monsters matching the query with XP up to maxXP.
func (s *Service) SearchMonsters(query string, maxXP int) []monster.Monster {
	return s.repo.Search(query, maxXP)
}

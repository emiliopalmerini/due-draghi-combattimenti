package monster

import (
	"testing"

	domain "github.com/emiliopalmerini/due-draghi-combattimenti/internal/domain/monster"
)

type mockRepo struct {
	monsters []domain.Monster
}

func (r *mockRepo) FindByMaxXP(maxXP int) []domain.Monster {
	var result []domain.Monster
	for _, m := range r.monsters {
		if m.XP <= maxXP {
			result = append(result, m)
		}
	}
	return result
}

func (r *mockRepo) Search(query string, maxXP int) []domain.Monster {
	var result []domain.Monster
	for _, m := range r.monsters {
		if m.XP <= maxXP {
			result = append(result, m)
		}
	}
	return result
}

func newTestService() *Service {
	repo := &mockRepo{
		monsters: []domain.Monster{
			{ID: "goblin", Name: "Goblin", CR: "1/4", XP: 50},
			{ID: "orc", Name: "Orco", CR: "1/2", XP: 100},
			{ID: "dragon", Name: "Drago Rosso Adulto", CR: "17", XP: 18000},
		},
	}
	return NewService(repo)
}

func TestSearchMonsters(t *testing.T) {
	svc := newTestService()

	results := svc.SearchMonsters("", 1_000_000)
	if len(results) != 3 {
		t.Errorf("expected 3 monsters, got %d", len(results))
	}

	results = svc.SearchMonsters("", 100)
	if len(results) != 2 {
		t.Errorf("expected 2 monsters with XP<=100, got %d", len(results))
	}
}

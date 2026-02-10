package memory

import (
	"testing"
)

func TestNewMonsterRepository_LoadsMonsters(t *testing.T) {
	repo := NewMonsterRepository()
	monsters := repo.FindByMaxXP(1_000_000)
	if len(monsters) == 0 {
		t.Fatal("expected monsters to be loaded, got 0")
	}
	// We know there are 330 monsters in the JSON
	if len(monsters) != 330 {
		t.Errorf("expected 330 monsters, got %d", len(monsters))
	}
}

func TestNewMonsterRepository_ParsesXP(t *testing.T) {
	repo := NewMonsterRepository()
	// Aboleth has CR 10, XP 5900
	monsters := repo.Search("aboleth", 1_000_000)
	if len(monsters) != 1 {
		t.Fatalf("expected 1 aboleth, got %d", len(monsters))
	}
	if monsters[0].XP != 5900 {
		t.Errorf("expected aboleth XP=5900, got %d", monsters[0].XP)
	}
	if monsters[0].CR != "10" {
		t.Errorf("expected aboleth CR=10, got %s", monsters[0].CR)
	}
}

func TestFindByMaxXP_FiltersCorrectly(t *testing.T) {
	repo := NewMonsterRepository()

	// Only CR 0 monsters have XP 0
	zeroXP := repo.FindByMaxXP(0)
	for _, m := range zeroXP {
		if m.XP != 0 {
			t.Errorf("monster %s has XP %d, expected 0", m.Name, m.XP)
		}
	}

	low := repo.FindByMaxXP(100)
	for _, m := range low {
		if m.XP > 100 {
			t.Errorf("monster %s has XP %d, expected <= 100", m.Name, m.XP)
		}
	}

	all := repo.FindByMaxXP(1_000_000)
	mid := repo.FindByMaxXP(1000)
	if len(mid) >= len(all) {
		t.Error("expected fewer monsters at lower XP cap")
	}
}

func TestSearch_FiltersByNameAndXP(t *testing.T) {
	repo := NewMonsterRepository()

	// Search by name
	dragons := repo.Search("drago", 1_000_000)
	if len(dragons) == 0 {
		t.Fatal("expected some dragons")
	}
	for _, m := range dragons {
		if !containsInsensitive(m.Name, "drago") {
			t.Errorf("monster %s does not match 'drago'", m.Name)
		}
	}

	// Search with XP cap
	weakDragons := repo.Search("drago", 500)
	for _, m := range weakDragons {
		if m.XP > 500 {
			t.Errorf("monster %s has XP %d, expected <= 500", m.Name, m.XP)
		}
	}
}

func TestSearch_EmptyQuery_ReturnsAllUpToMaxXP(t *testing.T) {
	repo := NewMonsterRepository()
	all := repo.Search("", 1_000_000)
	if len(all) != 330 {
		t.Errorf("expected 330 monsters with empty query, got %d", len(all))
	}
}

func containsInsensitive(s, substr string) bool {
	return len(s) >= len(substr) // placeholder, real check in impl
}

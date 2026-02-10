package memory

import (
	"embed"
	"encoding/json"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/emiliopalmerini/due-draghi-combattimenti/internal/domain/monster"
)

//go:embed data/monsters.json
var monstersFS embed.FS

var xpRegex = regexp.MustCompile(`PE\s+([\d.]+)`)

type jsonMonster struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Size     string `json:"size"`
	CR       string `json:"cr"`
	CRDetail string `json:"cr_detail"`
	AC       string `json:"ac"`
	HP       string `json:"hp"`
}

// MonsterRepository provides in-memory access to monster data.
type MonsterRepository struct {
	monsters []monster.Monster
}

// NewMonsterRepository loads monsters from embedded JSON.
func NewMonsterRepository() *MonsterRepository {
	data, err := monstersFS.ReadFile("data/monsters.json")
	if err != nil {
		log.Fatalf("failed to read embedded monsters.json: %v", err)
	}

	var raw []jsonMonster
	if err := json.Unmarshal(data, &raw); err != nil {
		log.Fatalf("failed to parse monsters.json: %v", err)
	}

	monsters := make([]monster.Monster, len(raw))
	for i, m := range raw {
		monsters[i] = monster.Monster{
			ID:   m.ID,
			Name: m.Name,
			Type: m.Type,
			Size: m.Size,
			CR:   m.CR,
			XP:   parseXP(m.CRDetail),
			AC:   m.AC,
			HP:   m.HP,
		}
	}

	sort.Slice(monsters, func(i, j int) bool {
		return monsters[i].XP < monsters[j].XP
	})

	return &MonsterRepository{monsters: monsters}
}

// parseXP extracts the XP value from a cr_detail string like "10 (PE 5.900; BC +4)".
// Italian uses dots as thousands separators, so "5.900" → 5900.
func parseXP(crDetail string) int {
	matches := xpRegex.FindStringSubmatch(crDetail)
	if len(matches) < 2 {
		return 0
	}
	// Remove dots (Italian thousands separator)
	numStr := strings.ReplaceAll(matches[1], ".", "")
	xp, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}
	return xp
}

func (r *MonsterRepository) FindByMaxXP(maxXP int) []monster.Monster {
	var result []monster.Monster
	for _, m := range r.monsters {
		if m.XP <= maxXP {
			result = append(result, m)
		}
	}
	return result
}

func (r *MonsterRepository) Search(query string, maxXP int) []monster.Monster {
	query = strings.ToLower(query)
	var result []monster.Monster
	for _, m := range r.monsters {
		if m.XP > maxXP {
			continue
		}
		if query == "" || strings.Contains(strings.ToLower(m.Name), query) {
			result = append(result, m)
		}
	}
	return result
}

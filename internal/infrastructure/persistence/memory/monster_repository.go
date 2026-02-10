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

type jsonAbilityScores struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

type jsonSavingThrows struct {
	Strength     string `json:"strength"`
	Dexterity    string `json:"dexterity"`
	Constitution string `json:"constitution"`
	Intelligence string `json:"intelligence"`
	Wisdom       string `json:"wisdom"`
	Charisma     string `json:"charisma"`
}

type jsonNamedDescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type jsonMonster struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Group               string                 `json:"group"`
	Type                string                 `json:"type"`
	Subtype             string                 `json:"subtype"`
	Size                string                 `json:"size"`
	Alignment           string                 `json:"alignment"`
	AC                  string                 `json:"ac"`
	Initiative          string                 `json:"initiative"`
	HP                  string                 `json:"hp"`
	Speed               string                 `json:"speed"`
	AbilityScores       jsonAbilityScores      `json:"ability_scores"`
	AbilityMods         jsonAbilityScores      `json:"ability_mods"`
	SavingThrows        jsonSavingThrows       `json:"saving_throws"`
	Skills              string                 `json:"skills"`
	Resistances         string                 `json:"resistances"`
	DamageImmunities    string                 `json:"damage_immunities"`
	ConditionImmunities string                 `json:"condition_immunities"`
	Senses              string                 `json:"senses"`
	Languages           string                 `json:"languages"`
	CR                  string                 `json:"cr"`
	CRDetail            string                 `json:"cr_detail"`
	Equipment           string                 `json:"equipment"`
	Traits              []jsonNamedDescription  `json:"traits"`
	Actions             []jsonNamedDescription  `json:"actions"`
	BonusActions        []jsonNamedDescription  `json:"bonus_actions"`
	Reactions           []jsonNamedDescription  `json:"reactions"`
	LegendaryActions    []jsonNamedDescription  `json:"legendary_actions"`
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

			Group:               m.Group,
			Subtype:             m.Subtype,
			Alignment:           m.Alignment,
			Initiative:          m.Initiative,
			Speed:               m.Speed,
			AbilityScores:       monster.AbilityScores(m.AbilityScores),
			AbilityMods:         monster.AbilityScores(m.AbilityMods),
			SavingThrows:        monster.SavingThrows(m.SavingThrows),
			Skills:              m.Skills,
			Senses:              m.Senses,
			Languages:           m.Languages,
			Resistances:         m.Resistances,
			DamageImmunities:    m.DamageImmunities,
			ConditionImmunities: m.ConditionImmunities,
			Equipment:           m.Equipment,
			CRDetail:            m.CRDetail,
			Traits:              convertNamedDescriptions(m.Traits),
			Actions:             convertNamedDescriptions(m.Actions),
			BonusActions:        convertNamedDescriptions(m.BonusActions),
			Reactions:           convertNamedDescriptions(m.Reactions),
			LegendaryActions:    convertNamedDescriptions(m.LegendaryActions),
		}
	}

	sort.Slice(monsters, func(i, j int) bool {
		return monsters[i].XP < monsters[j].XP
	})

	return &MonsterRepository{monsters: monsters}
}

func convertNamedDescriptions(src []jsonNamedDescription) []monster.NamedDescription {
	if len(src) == 0 {
		return nil
	}
	result := make([]monster.NamedDescription, len(src))
	for i, nd := range src {
		result[i] = monster.NamedDescription{Name: nd.Name, Description: nd.Description}
	}
	return result
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

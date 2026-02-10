package monster

// Monster represents a D&D creature with combat-relevant stats.
type Monster struct {
	ID   string
	Name string
	Type string
	Size string
	CR   string
	XP   int
	AC   string
	HP   string
}

// Repository defines the interface for accessing monster data.
type Repository interface {
	FindByMaxXP(maxXP int) []Monster
	Search(query string, maxXP int) []Monster
}

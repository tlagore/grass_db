package db_manager

type GrassEntry struct {
	GenusSpecies string
	IsPerennial bool
	CulmDensity string
	RootingCharactersitic string
	CulmGrowth string
	CulmLengthMinCm int
	CulmLengthMaxCm int
	CulmDiameterMinMm int
	CulmDiameterMaxMm int
	IsWoody bool
	CulmInternode string
	LocationBroad string
	LocationNarrow string
	Notes string
}
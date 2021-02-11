package db_manager

type GrassEntry struct {
	GrassSpecies          string
	IsPerennial           bool
	IsAnnual              bool
	CulmDensity           string
	RootingCharactersitic string
	CulmGrowth            string
	CulmLengthMinCm       float64
	CulmLengthMaxCm       float64
	CulmDiameterMinMm     float64
	CulmDiameterMaxMm     float64
	IsWoody               bool
	CulmInternode 		  string
	LocationBroad 		  string
	LocationNarrow 		  string
	Notes 				  string
}

type BambooEntry struct {
	GenusSpecies     		string
	IsInvasive       		bool
	DisputedNativeRange 	bool
	NumIntroductions 		int
}
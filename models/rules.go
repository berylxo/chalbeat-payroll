package models

type PayeBand struct {
	Min  float64
	Max  float64
	Rate float64
}

type Rules struct {
	Paye struct {
		Bands          []PayeBand
		PersonalRelief float64
	}

	Nssf struct {
		Tier1Limit float64
		Tier1Rate  float64
		Tier2Rate  float64
		Tier2Cap   float64
	}

	HousingLevyRate float64
	ShifRate        float64
}

func DefaultRules() Rules {
	return Rules{
		HousingLevyRate: 0.015,
		ShifRate:        0.0275,
		Paye: struct {
			Bands          []PayeBand
			PersonalRelief float64
		}{
			Bands: []PayeBand{
				{Min: 0, Max: 24000, Rate: 0.1},
				{Min: 24001, Max: 32333, Rate: 0.25},
				{Min: 32334, Max: 9999999, Rate: 0.3},
			},
			PersonalRelief: 2400,
		},
		Nssf: struct {
			Tier1Limit float64
			Tier1Rate  float64
			Tier2Rate  float64
			Tier2Cap   float64
		}{
			Tier1Limit: 8000,
			Tier1Rate:  0.06,
			Tier2Rate:  0.06,
			Tier2Cap:   72000,
		},
	}
}

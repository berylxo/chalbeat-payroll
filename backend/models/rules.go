package models

type PayeBand struct {
	Width int64   // The maximum amount of money that can fit in this band
	Rate  float64 // The tax rate applied to this band
}

type Rules struct {
	HousingLevyRate float64
	ShifRate        float64
	PersonalRelief  int64
	NssfUpperLimit  int64
	PayeBands       []PayeBand
}

func DefaultRules() Rules {
	return Rules{
		HousingLevyRate: 0.015,
		ShifRate:        0.0275,
		PersonalRelief:  240000,
		NssfUpperLimit:  7200000,
		PayeBands: []PayeBand{
			{Width: 2400000, Rate: 0.10},   // First 24,000
			{Width: 833300, Rate: 0.25},    // Next 8,333 (up to 32,333)
			{Width: 46766700, Rate: 0.30},  // Next 467,667 (up to 500,000)
			{Width: 30000000, Rate: 0.325}, // Next 300,000 (up to 800,000)
			{Width: -1, Rate: 0.35},        // -1 signifies infinity/taxable balance
		},
	}
}

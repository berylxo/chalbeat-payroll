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

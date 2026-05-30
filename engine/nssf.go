package engine

import "github.com/berylxo/chalbeat-payroll/models"

func CalculateNSSF(gross float64, rules models.Rules) float64 {
	tier1 := min(gross, rules.Nssf.Tier1Limit) * rules.Nssf.Tier1Rate

	tier2 := 0.0
	if gross > rules.Nssf.Tier1Limit {
		tier2Base := min(gross, rules.Nssf.Tier2Cap) - rules.Nssf.Tier1Limit
		tier2 = tier2Base * rules.Nssf.Tier2Rate
	}

	return tier1 + tier2
}

package engine

import "github.com/berylxo/chalbeat-payroll/models"

func CalculateNSSF(gross int64, rules models.Rules) int64 {
	// 6% of gross income up to the upper limit of KES 72,000
	pensionable := gross
	if pensionable > rules.NssfUpperLimit {
		pensionable = rules.NssfUpperLimit
	}
	return calculatePercentage(pensionable, 0.06)
}

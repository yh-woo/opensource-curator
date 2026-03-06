package scoring

// ScoreResult holds the complete scoring output for a library.
type ScoreResult struct {
	Overall           float64 `json:"overall"`
	MaintenanceHealth float64 `json:"maintenanceHealth"`
	APIClarity        float64 `json:"apiClarity"`
	DocQuality        float64 `json:"docQuality"`
	SecurityPosture   float64 `json:"securityPosture"`
	CommunitySignal   float64 `json:"communitySignal"`
	DeprecationSafety float64 `json:"deprecationSafety"`
	Deprecated        bool    `json:"deprecated"`
	DeprecationReason string  `json:"deprecationReason,omitempty"`
	Version           string  `json:"version"`
}

// Calculate computes the full score for a library from collected raw data.
// data should be the merged output from all collectors (GitHub + npm).
func Calculate(data map[string]any) ScoreResult {
	// Step 1: Pre-filter
	pf := PreFilter(data)
	if pf.Deprecated {
		return ScoreResult{
			Overall:           0,
			Deprecated:        true,
			DeprecationReason: pf.Reason,
			Version:           Version,
		}
	}

	// Step 2: Calculate individual metrics
	maintenance := CalcMaintenanceHealth(data)
	apiClarity := CalcAPIClarity(data)
	docQuality := CalcDocQuality(data)
	security := CalcSecurityPosture(data)
	community := CalcCommunitySignal(data)
	deprecation := CalcDeprecationSafety(data)

	// Step 3: Weighted average
	overall := maintenance*WeightMaintenanceHealth +
		apiClarity*WeightAPIClarity +
		docQuality*WeightDocQuality +
		security*WeightSecurityPosture +
		community*WeightCommunitySignal +
		deprecation*WeightDeprecationSafety

	return ScoreResult{
		Overall:           round2(overall),
		MaintenanceHealth: round2(maintenance),
		APIClarity:        round2(apiClarity),
		DocQuality:        round2(docQuality),
		SecurityPosture:   round2(security),
		CommunitySignal:   round2(community),
		DeprecationSafety: round2(deprecation),
		Deprecated:        false,
		Version:           Version,
	}
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

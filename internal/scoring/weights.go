package scoring

const Version = "1.0.0"

// Weights for the scoring algorithm (must sum to 1.0).
const (
	WeightMaintenanceHealth = 0.25
	WeightAPIClarity        = 0.20
	WeightDocQuality        = 0.15
	WeightSecurityPosture   = 0.15
	WeightCommunitySignal   = 0.15
	WeightDeprecationSafety = 0.10
)

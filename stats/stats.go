package stats

// PendingStat is a pending Stat struct
type PendingStat struct {
	ValidShares      uint64
	StaleShares      uint64
	InvalidShares    uint64
	ReportedHashrate float64
	IPAddress        string
}

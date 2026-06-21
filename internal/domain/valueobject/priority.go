package valueobject

type Priority string

const (
	PriorityLow      Priority = "LOW"
	PriorityMedium   Priority = "MEDIUM"
	PriorityHigh     Priority = "HIGH"
	PriorityCritical Priority = "CRITICAL"
)

func (p Priority) IsValid() bool {
	return p == PriorityLow ||
		p == PriorityMedium ||
		p == PriorityHigh ||
		p == PriorityCritical
}

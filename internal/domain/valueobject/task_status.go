package valueobject

type TaskStatus string

const (
	StatusTodo       TaskStatus = "TODO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusReview     TaskStatus = "REVIEW"
	StatusDone       TaskStatus = "DONE"
)

func (s TaskStatus) IsValid() bool {
	return s == StatusTodo ||
		s == StatusInProgress ||
		s == StatusReview ||
		s == StatusDone
}

func (s TaskStatus) CanTransitionTo(next TaskStatus) bool {
	switch s {
	case StatusTodo:
		return next == StatusInProgress
	case StatusInProgress:
		return next == StatusReview
	case StatusReview:
		return next == StatusDone || next == StatusInProgress
	default:
		return false
	}
}

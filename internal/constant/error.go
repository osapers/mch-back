package constant

// ErrorNoResult describes error when db query returns no results
type ErrorNoResult struct {
	Target string
}

func (e *ErrorNoResult) Error() string {
	return "no results"
}

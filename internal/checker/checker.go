package checker

// CheckRunner provides actions necessary to support
// running various types of checks.
type CheckRunner interface {
	PerformCheck() error
}

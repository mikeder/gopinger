package checker

// CheckRunner actions necessary to support running various types of checks.
type CheckRunner interface {
	PerformCheck(check interface{}) (interface{}, error)
}

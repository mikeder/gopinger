package checks


type Check struct {
  Url           string
  HealthyCode   int
  Method        string
  Timeout       int
}

type Result struct {
  Status        string
  Reason        string
}

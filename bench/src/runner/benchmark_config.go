package runner

type BenchmarkConfig struct {
	BaseURL               string
	ScenarioName          string
	CPUCount              int
	InitialRate           int
	DurationSeconds       int
	RequestTimeoutSeconds int
}

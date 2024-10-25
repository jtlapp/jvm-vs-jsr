package runner

type BenchmarkConfig struct {
	BaseURL               string
	ScenarioName          string
	CPUCount              int
	MaxConnections        int
	InitialRate           int
	DurationSeconds       int
	RequestTimeoutSeconds int
	MinWaitSeconds        int
}

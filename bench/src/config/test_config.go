package config

type TestConfig struct {
	ScenarioName             string
	CPUsToUse                int
	WorkerCount              int
	MaxConnections           int
	InitialRequestsPerSecond int
	DurationSeconds          int
	RequestTimeoutSeconds    int
	MinWaitSeconds           int
}

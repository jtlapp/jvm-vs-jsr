package config

type TestConfig struct {
	ScenarioName             string
	CPUsToUse                int
	WorkerCount              int
	MaxConnections           int
	InitialRequestsPerSecond int
	InitialRandomSeed        int
	DurationSeconds          int
	RequestTimeoutSeconds    int
	MinWaitSeconds           int
}

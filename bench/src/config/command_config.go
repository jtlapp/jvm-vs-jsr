package config

type CommandConfig struct {
	ConfigFile               *string
	ScenarioName             *string
	LongSleepMillis          *int
	ShortSleepMillis         *int
	PercentLongRequests      *float64
	CPUsToUse                *int
	WorkerCount              *int
	MaxConnections           *int
	InitialRequestsPerSecond *int
	InitialRandomSeed        *int
	DurationSeconds          *int
	RequestTimeoutSeconds    *int
	MinSecondsBetweenTests   *int
	LoopCount                *int
	ResetRandomSeed          *bool
	TrialCount               *int
}

package config

type CommandConfig struct {
	ConfigFile               *string
	ScenarioName             *string
	LongSleepMillis          *int
	ShortSleepMillis         *int
	PercentLongRequests      *int
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

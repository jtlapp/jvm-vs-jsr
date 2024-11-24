package config

type ScenarioConfig struct {
	LongSleepMillis     int
	ShortSleepMillis    int
	PercentLongRequests int
}

func NewScenarioConfig(commandConfig *CommandConfig) *ScenarioConfig {
	return &ScenarioConfig{
		LongSleepMillis:     *commandConfig.LongSleepMillis,
		ShortSleepMillis:    *commandConfig.ShortSleepMillis,
		PercentLongRequests: *commandConfig.PercentLongRequests,
	}
}

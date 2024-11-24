package config

import "jvm-vs-jsr.jtlapp.com/benchmark/command/usage"

type ScenarioConfig struct {
	LongSleepMillis     int
	ShortSleepMillis    int
	PercentLongRequests int
}

func NewScenarioConfig(commandConfig *usage.CommandConfig) *ScenarioConfig {
	return &ScenarioConfig{
		LongSleepMillis:     *commandConfig.LongSleepMillis,
		ShortSleepMillis:    *commandConfig.ShortSleepMillis,
		PercentLongRequests: *commandConfig.PercentLongRequests,
	}
}

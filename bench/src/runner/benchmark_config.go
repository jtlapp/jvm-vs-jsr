package runner

type BenchmarkConfig struct {
	ClientVersion		  string
	BaseAppUrl            string
	AppName			   string
	AppVersion			string
	AppConfig			 map[string]interface{}
	ScenarioName          string
	CPUCount              int
	MaxConnections        int
	InitialRate           int
	DurationSeconds       int
	RequestTimeoutSeconds int
	MinWaitSeconds        int
}

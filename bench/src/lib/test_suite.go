package lib

import vegeta "github.com/tsenart/vegeta/lib"

type TestSuite interface {
	GetName() string
	Init(backendDB *BackendDB) error
	SetUpTestTables() error
	SetSharedQueries() error
	GetTargetProvider(baseUrl string) func(*vegeta.Target) error
}

type TestSuiteFactory func() (TestSuite, error)

package lib

import vegeta "github.com/tsenart/vegeta/lib"

type TestSuite interface {
	GetName() string
	Init() error
	SetUpDatabase() error
	SetSharedQueries() error
	GetTargeter(baseUrl string) vegeta.Targeter
}

type TestSuiteFactory func() (TestSuite, error)

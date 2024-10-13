package orderitems

import (
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-js.jtlapp.com/benchmark/lib"
)

type Suite struct {
	databaseSetup lib.DatabaseSetup
}

func (s *Suite) GetName() string {
	return "orderitems"
}

func (s *Suite) Init() error {
	impl := &SetupImpl{}

	databaseSetup, err := lib.NewDatabaseSetup(impl)
	if err != nil {
		return err
	}
	s.databaseSetup = *databaseSetup
	return nil
}

func (s *Suite) SetUpDatabase() error {
	return s.databaseSetup.PopulateDatabase()
}

func (s *Suite) SetSharedQueries() error {
	return s.databaseSetup.CreateSharedQueries()
}

func (s *Suite) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	test := NewBenchmarkTest(baseUrl)
	return func(target *vegeta.Target) error {
		*target = *test.getRequest()
		return nil
	}
}

func (s *Suite) Close() error {
	return s.databaseSetup.Release()
}

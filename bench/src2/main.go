package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	taggedints "bench.bin/tagged-ints"
	// Import test suites
)

// Options struct to hold the test options
type Options struct {
	verbose bool
	retry   int
}

// loadOptions parses the test options
func loadOptions() Options {
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	retry := flag.Int("retry", 1, "Number of retries")
	flag.Parse()
	return Options{
		verbose: *verbose,
		retry:   *retry,
	}
}

// TestSuite interface that every test suite must implement
type TestSuite interface {
	Name() string
	CreateSetup()
	Test(options Options)
}

// testSuitesSlice holds pointers to all test suites
var testSuitesSlice = []TestSuite{
	&taggedints.Suite{},
}

// getTestSuite finds the test suite by name in the slice
func getTestSuite(name string) (TestSuite, bool) {
	for _, suite := range testSuitesSlice {
		if suite.Name() == name {
			return suite, true
		}
	}
	return nil, false
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <test-suite-name> setup|test [<options>]", os.Args[0])
	}

	// Get the test suite name and the mode (setup or test)
	suiteName := os.Args[1]
	mode := os.Args[2]

	// Find the test suite from the slice
	suite, valid := getTestSuite(suiteName)
	if !valid {
		log.Fatalf("Unknown test suite: %s", suiteName)
	}

	// Handle the "setup" and "test" modes
	if mode == "test" {
		// Load test options if in test mode
		options := loadOptions()
		suite.Test(options)
	} else if mode == "setup" {
		suite.Setup()
	} else {
		log.Fatalf("Invalid argument '%s'. Must be 'setup' or 'test'.", mode)
	}
}


func setup() {
	const DATABASE_URL = "postgres://pgbouncer-service:6432/testdb"
	const USERNAME = "user"
	const PASSWORD = "password"

	setup, err := CreateSetup(DATABASE_URL, USERNAME, PASSWORD)
	if err != nil {
		log.Fatalf("Failed to create setup: %v", err)
	}

	err = setup.Run()
	if err != nil {
		log.Fatalf("Failed to run setup: %v", err)
	}

	err = setup.Release()
	if err != nil {
		log.Fatalf("Failed to release setup: %v", err)
	}

	fmt.Printf("'%s' database setup completed.\n", setup.GetName())
}

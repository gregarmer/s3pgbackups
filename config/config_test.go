package config

import "testing"

func TestPreFlight(t *testing.T) {
	conf := Config{}
	err := conf.PreFlight()
	if err == nil {
		t.Fatalf("error should be set")
	}
}

func TestCopy(t *testing.T) {
	conf := Config{S3Bucket: "foo"}
	c := conf.Copy()
	c.S3Bucket = "bar"
	if conf.S3Bucket == c.S3Bucket {
		t.Fatalf("config.Copy() should return a copy that doesn't affect the original")
	}
}

func TestShouldExcludeDbWithNoExcludes(t *testing.T) {
	conf := Config{}
	var testDatabase string

	// includes
	testDatabase = "foobar"
	if conf.ShouldExcludeDb(testDatabase) != false {
		t.Fatalf("%s should not have been excluded.", testDatabase)
	}
}

func TestShouldExcludeDb(t *testing.T) {
	conf := Config{Excludes: []string{"foo", "*.bar", "baz.*", "quux.*"}}
	var testDatabase string

	// excludes
	testDatabase = "foo"
	if conf.ShouldExcludeDb(testDatabase) != true {
		t.Fatalf("%s should have been excluded.", testDatabase)
	}

	testDatabase = "quux"
	if conf.ShouldExcludeDb(testDatabase) != true {
		t.Fatalf("%s should have been excluded.", testDatabase)
	}

	// includes
	testDatabase = "foobar"
	if conf.ShouldExcludeDb(testDatabase) != false {
		t.Fatalf("%s should not have been excluded.", testDatabase)
	}
}

func TestShouldExcludeTable(t *testing.T) {
	conf := Config{Excludes: []string{"foo", "*.bar", "baz.*", "quux.gah"}}
	var testDatabase string
	var testTable string

	// excludes
	testDatabase = "foo"
	testTable = "bar"
	if conf.ShouldExcludeTable(testDatabase, testTable) != true {
		t.Fatalf("%s.%s should have been excluded.", testDatabase, testTable)
	}

	testDatabase = "baz"
	testTable = "quux"
	if conf.ShouldExcludeTable(testDatabase, testTable) != true {
		t.Fatalf("%s.%s should have been excluded.", testDatabase, testTable)
	}

	testDatabase = "quux"
	testTable = "gah"
	if conf.ShouldExcludeTable(testDatabase, testTable) != true {
		t.Fatalf("%s.%s should have been excluded.", testDatabase, testTable)
	}

	// includes
	testDatabase = "foo"
	testTable = "gah"
	if conf.ShouldExcludeTable(testDatabase, testTable) != false {
		t.Fatalf("%s.%s should not have been excluded.", testDatabase, testTable)
	}

	testDatabase = "foobar"
	testTable = "gah"
	if conf.ShouldExcludeTable(testDatabase, testTable) != false {
		t.Fatalf("%s.%s should not have been excluded.", testDatabase, testTable)
	}
}

package utils

import (
	"log"
	"os"
)

func Fatalf(error string, args ...interface{}) {
	criticalLog := log.New(os.Stderr, "", log.LstdFlags)
	criticalLog.Printf(error, args...)
	os.Exit(1)
}

func CheckErr(err error) {
	if err != nil {
		Fatalf("Error: %s", err)
	}
}

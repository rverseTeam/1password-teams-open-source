package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

func printUsageAndExit() {
	log.Fatal("Usage: ./processor <review|approve> [--test <name>]")
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("missing variable %s", key)
	}

	return value, nil
}

func main() {
	sentryDsn, err := getEnv("SENTRY_DSN")
	if err != nil {
		log.Fatal(err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn: sentryDsn,
	})
	if err != nil {
		log.Fatalf("Sentry initialization failed: %v\n", err)
	}

	defer sentry.Flush(2 * time.Second)

	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	command := os.Args[1]

	if len(os.Args) == 4 && os.Args[2] == "--test" {
		setTestApplication(os.Args[3])
	} else if len(os.Args) != 2 {
		printUsageAndExit()
	}

	github := GitHub{}
	application := Application{}

	switch command {
	case "review":
		reviewer := Reviewer{
			gitHub:      &github,
			application: &application,
		}
		reviewer.Review()
	case "approve":
		approver := Approver{
			gitHub:      &github,
			application: &application,
		}
		approver.Approve()
	default:
		fmt.Printf("Invalid command: %s\n", command)
		printUsageAndExit()
	}
}

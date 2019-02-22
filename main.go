package main

import (
	a "./awshelpers"
	t "./tablehelpers"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

// Lambda handler function
// Connects to the RDS instance and the database mentioned
// Fetches the tables that are to be checked during this hour
// Checks the diff of the two
// If there is a mismatch, publish to the given SNS channel
func monitorTables(ctx context.Context) {
	log.Printf("Starting execution in %s :: %s", os.Getenv("AWS_REGION"), os.Getenv("ENVIRONMENT"))

	scheduledTables := t.FetchScheduledTables()
	log.Print("Tables checked :: ", scheduledTables)

	if len(scheduledTables) > 0 {
		dbTables := a.FetchTablesFromDB()

		failedTables := t.SliceDifference(scheduledTables, dbTables)
		log.Print("Tables in DB :: ", len(dbTables))
		log.Print("Failed table creations :: ", failedTables)

		if len(failedTables) > 0 {
			a.PublishResultsToSns(failedTables)
		}
	}
}

// Invokes the lambda handler
func main() {
	lambda.Start(monitorTables)
}

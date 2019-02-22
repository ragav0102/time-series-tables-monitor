package main

import (
	a "./awshelpers"
	t "./tablehelpers"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

// Event must contain all the fields
type TriggerEvent struct {
	Name        string   `json:"name"`
	ServiceName string   `json:"service_name"`
	Env         string   `json:"env"`
	Cron        string   `json:"cron"`
	ModelNames  []string `json:"model_names"`
	RdsCreds    string   `json:"rds_creds"`
	SnsTopicArn string   `json:"sns_topic_arn"`
}

// Lambda handler function
// Connects to the RDS instance and the database mentioned
// Fetches the tables that are to be checked during this hour
// Checks the diff of the two
// If there is a mismatch, publish to the given SNS channel
func monitorTables(ctx context.Context, params TriggerEvent) {
	log.Printf("Starting execution in %s :: %s", os.Getenv("AWS_REGION"), params.Env)

	scheduledTables := t.FetchScheduledTables(params.ModelNames, params.Cron)
	log.Print("Tables checked :: ", scheduledTables)

	if len(scheduledTables) > 0 {
		dbTables := a.FetchTablesFromDB(params.RdsCreds)

		failedTables := t.SliceDifference(scheduledTables, dbTables)
		log.Print("Tables in DB :: ", len(dbTables))
		log.Print("Failed table creations :: ", failedTables)

		if len(failedTables) > 0 {
			a.PublishResultsToSns(params.Name, failedTables, params.Env, params.SnsTopicArn, params.ServiceName)
		}
	}
}

// Invokes the lambda handler
func main() {
	lambda.Start(monitorTables)
}

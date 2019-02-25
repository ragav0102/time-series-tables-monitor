package awshelpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"time"
)

// Connects to the given DB and returns all the tables present
func FetchTablesFromDB(dbCreds string) []string {
	var dbTables []string
	// dbCreds := os.Getenv("RDS_CREDS")
	db, err := sql.Open("mysql", dbCreds)

	if err != nil {
		log.Print(err.Error())
		return dbTables
	} else {
		log.Print("Connected to DB")
	}

	dbName := strings.Split(dbCreds, "/")[1]
	log.Print("Checking DB :: ", dbName)
	sqlQuery := fmt.Sprintf("select table_name from information_schema.tables where table_schema='%s'", dbName)
	rows, err := db.Query(sqlQuery)

	defer rows.Close()
	for rows.Next() {
		var (
			table_name string
		)
		if err := rows.Scan(&table_name); err != nil {
			log.Fatal(err)
		}
		dbTables = append(dbTables, table_name)
	}
	return dbTables
}

func PublishResultsToSns(event string, tables []string, env string, topicArn string, serviceName string) {
	snsSvc := sns.New(session.New())
	message, _ := json.Marshal(map[string]interface{}{
		"event":          event,
		"service_name":   serviceName,
		"missing_tables": tables,
		"time":           time.Now().Format("2006-01-02T15:04:05.999999"),
		"region":         os.Getenv("AWS_REGION"),
		"environment":    env,
	})

	publishParams := &sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(topicArn),
	}

	resp, err := snsSvc.Publish(publishParams)

	if err != nil {
		log.Print("Error during notifying :: ", err.Error())
		return
	}

	log.Print("Successfully published to topic :: ", resp)
}

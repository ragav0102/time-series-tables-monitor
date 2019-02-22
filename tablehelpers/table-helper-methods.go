package tablehelpers

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Gives the suffix for the weekly tables that will be used during the successive week
func weeklyTableSuffix() string {
	nextWeek := time.Now().UTC().AddDate(0, 0, 7)
	return nextWeek.AddDate(0, 0, (1 - int(nextWeek.Weekday()))).Format("2006_01_02")
}

// Gives the suffix for the daily tables that will be used during the next day
func dailyTableSuffix() string {
	nextDay := time.Now().UTC().AddDate(0, 0, 1)
	return nextDay.Format("2006_01_02")
}

// Gives the suffix for the hourly tables that will be used during the next hour
func hourlyTableSuffix() string {
	nextHour := time.Now().UTC().Add(1 * time.Hour)
	h := strconv.Itoa(nextHour.Hour())
	return nextHour.Format("2006_01_02") + "_" + h
}

func fetchTableTimeSuffix(tableType string) string {
	if tableType == "hourly" {
		return "_" + hourlyTableSuffix()
	} else if tableType == "daily" {
		return "_" + dailyTableSuffix()
	} else {
		return "_" + weeklyTableSuffix()
	}
}

func isHourly(m int, cronExp []int) bool {
	return cronExp[0] != -1 && m >= cronExp[0]
}

func isDaily(m int, h int, cronExp []int) bool {
	if cronExp[1] == -1 {
		return false
	}
	if cronExp[0] == -1 {
		return h >= cronExp[1]
	} else {
		return m >= cronExp[0] && h >= cronExp[1]
	}
}

func isWeekly(m int, h int, wd int, cronExp []int) bool {
	if cronExp[2] == -1 {
		return false
	}
	if cronExp[1] == -1 {
		return wd == cronExp[2]
	} else {
		return isDaily(m, h, cronExp) && wd == cronExp[2]
	}
}

// Parses the CRON expression given in format:
//     minute.hour.weekDay
//   and returns array of int values
// If any of the value is not an int,
//    -1 is will be the placeholder
func parsedCronExp() []int {
	var result []int
	cronExp := strings.Split(os.Getenv("CRON_M_H_WD"), ".")

	for _, i := range cronExp {
		j, err := strconv.Atoi(i)

		if err != nil {
			log.Print(err)
			result = append(result, -1)
		} else {
			result = append(result, j)
		}
	}
	return result
}

// Method to check missed tables
// When given two slices returns the missing items
//   from slice2 that is present in slice1
func SliceDifference(slice1 []string, slice2 []string) []string {
	var diff []string

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

// Checks CRON expression given
//   and returns the table types to be checked now
func fetchCheckableTableTypes() []string {
	var tableTypes []string

	cronExp := parsedCronExp()
	log.Print("Preparing tables based on CRON expression", cronExp)

	t := time.Now()
	m, h, wd := t.Minute(), t.Hour(), int(t.Weekday())

	// based on CRON time appends types of tables to be checked
	if isHourly(m, cronExp) {
		tableTypes = append(tableTypes, "hourly")
	}

	if isDaily(m, h, cronExp) {
		tableTypes = append(tableTypes, "daily")

	}

	if isWeekly(m, h, wd, cronExp) {
		tableTypes = append(tableTypes, "weekly")
	}

	return tableTypes
}

// Returns a slice of table names that will be checked
//   by lambda for current time
//   based on current time, CRON expression
//   and models present as an ENV variable
func FetchScheduledTables() []string {
	var tables []string

	tableTypes := fetchCheckableTableTypes()
	models := strings.Split(os.Getenv("MODEL_NAMES"), ",")

	for _, model := range models {
		for _, tableType := range tableTypes {
			tableName := model + "_" + tableType + fetchTableTimeSuffix(tableType)
			tables = append(tables, tableName)
		}
	}

	return tables
}

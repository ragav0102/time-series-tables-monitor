# time-series-tables-monitor
AWS Lambda handler for Golang, which will monitor table creations in a Time-series Mysql DB and notify through a AWS SNS topic, when there is a failure in table creation CRON.



**Runs on Go 1.x runtime**


#### Handler name: 'main'



### Required ENV variables:

  * **ENVIRONMENT**   - To identify environment while debugging

  * **SNS_TOPIC_ARN** - Topic ARN for SNS notification

  * **RDS_CREDS**     - Credentials for connecting to RDS instance.

    Format: "username:password@tcp(host:port)/db_name"

  * **MODEL_NAMES**   - Table prefixes that the function adds to timestamp as comma-separated values

    ex: For a daily table "group_2019_01_02", this value would be group

  * **CRON_M_H_WD**   - CRON expression of the job creating hourly, daily and/or weekly tables

    Format: "minute.hour.week_day"



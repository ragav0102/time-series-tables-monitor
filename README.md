# time-series-tables-monitor
AWS Lambda handler for Golang, which will monitor table creations in a Time-series Mysql DB and notify through a AWS SNS topic, when there is a failure in table creation CRON.



**Runs on Go 1.x runtime**


#### Handler name: 'main'



##### Lambda can be triggered through AWS Cloudwatch Rules based on our own requirement with a JSON payload


#### Sample payload

```
{
  "name": "Fire cannot kill a dragon",
  "service_name": "ragav's test",
  "env": "staging",
  "cron": "25.13.1",
  "model_names": ["group","people","targaryen"],
  "rds_creds": "userdas:passj32j13@tcp(dragonstone.targaryan.ds:2000)/db_namehdbyr47",
  "sns_topic_arn": "arn:aws:sns:us-east-1:7423463:viserion"
}
```

### Required fields:

  * **env**   - To identify environment while debugging

  * **sns_topic_arn** - Topic ARN for SNS notification

  * **rds_creds**     - Credentials for connecting to RDS instance.

    ``` "username:password@tcp(host:port)/db_name" ```

  * **model_names**   - Table prefixes that the function adds to timestamp as comma-separated values

    ex: For a daily table *group_2019_01_02*, this value would be group

  * **cron**   - CRON expression of the job creating hourly, daily and/or weekly tables

     ``` "minute.hour.week_day" ```



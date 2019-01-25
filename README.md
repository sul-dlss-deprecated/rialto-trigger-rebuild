# RIALTO Trigger Rebuild

The purpose of this program is to do a full rebuild of the RIALTO derivative data store(s) based on data from the canonical data store.

## Scheduled messaging

This ECS task can be scheduled to run on a recurring basis through [CloudWatch Events](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/WhatIsCloudWatchEvents.html).

Note: The `rebuildTrigger` event is currently disabled.

Resources: [Scheduled Event Rules](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html)

## Build

```
docker build -t suldlss/rialto-trigger-rebuild:latest .
```

## Deploy
```
docker push suldlss/rialto-trigger-rebuild:latest
```

## Add Schedule Event

See https://docs.aws.amazon.com/AmazonECS/latest/developerguide/scheduled_tasks.html

## ENV variables required for task

```
RIALTO_SNS_ENDPOINT=<ENDPOINT>
RIALTO_TOPIC_ARN=<ARN>
SPARQL_ENDPOINT=<SPARQL>
SOLR_HOST=<SOLR>
SOLR_COLLECTION=<COLLECTION>
RDS_USERNAME=<USERNAME>
RDS_PASSWORD=<PASSWORD>
RDS_DB_NAME=<DBNAME>
RDS_HOSTNAME=<HOST>
RDS_PORT=<PORT>
RIALTO_TRIPLELIMIT=<TRIPLELIMIT>
```

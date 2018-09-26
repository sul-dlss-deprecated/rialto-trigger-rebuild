# RIALTO Trigger Rebuild Lambda

The purpose of this lambda is to do a full rebuild of the RIALTO derivative data store(s) based on data from the canonical data store.

## Scheduled messaging

This lambda can be scheduled to run on a recurring basis through [CloudWatch Events](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/WhatIsCloudWatchEvents.html).

Note: The `rebuildTrigger` event is currently disabled.

Resources: [Scheduled Event Rules](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html)

## Manual messaging

An API Endpoint is available for this lambda in order trigger a rebuild on demand.

```
curl -H "X-Api-Key: <KEY>" -X POST https://4klo9rqst3.execute-api.us-west-2.amazonaws.com/development/rialto-trigger-rebuild
```

## Build Lambda

```
make
```

## Upload Lambda

```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws \
--endpoint-url http://localhost:4574 lambda create-function \
--function-name triggerRebuild \
--runtime go1.x \
--role RialtoLambda \
--handler main \
--environment "Variables={RIALTO_SNS_ENDPOINT=<ENDPOINT>,RIALTO_TOPIC_ARN=<ARN>, \
  SPARQL_ENDPOINT=<SPARQL>,SOLR_HOST=<SOLR>,SOLR_COLLECTION=<COLLECTION>,\
  RDS_USERNAME=<USERNAME>,RDS_PASSWORD=<PASSWORD>,RDS_DB_NAME=<DBNAME>,\
  RDS_HOSTNAME=<HOST>,RDS_PORT=<PORT>}" \
--zip-file fileb://lambda.zip
```

Note: When deploying to AWS proper, additional network and subnet settings are required. All of the above can also be manually accomplished through the console.

## Add Schedule Event

```
aws events put-rule --name "rebuildTrigger" --schedule-expression "cron(25 21 ? * * *)"
```

```
aws events put-targets --rule rebuildTrigger \
  --targets "Id"="1", \
  "Arn"="arn:aws:lambda:us-east-1:123456789012:function:triggerRebuild"
```

## ENV variables required for lambda

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
```

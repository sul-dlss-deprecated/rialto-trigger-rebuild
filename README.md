# RIALTO Trigger Rebuild Lambda

The purpose of this lambda is to publish a message to the configured SNS topic that will be used to trigger a full rebuild of the RIALTO derivative data store(s) based on data from the canonical data store.

## Scheduled messaging

This lambda can be scheduled to run on a recurring basis through [CloudWatch Events](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/WhatIsCloudWatchEvents.html).

Note: The `rebuildTrigger` event is currently disabled.

Resources: [Scheduled Event Rules](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html)

## Manual messaging

An API Endpoint is available for this lamba in order trigger a rebuild on demand. (link not included for security)

## Build Lambda

```
GOOS=linux go build -o main
zip lambda.zip main
```

## Upload Lambda

```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws \
--endpoint-url http://localhost:4574 lambda create-function \
--function-name triggerRebuild \
--runtime go1.x \
--role RialtoLambda \
--handler main \
--environment "Variables={REBUILD_ACTION=rebuild,REBUILD_MESSAGE:full,RIALTO_SNS_ENDPOINT=<ENDPOINT>,RIALTO_TOPIC_ARN=<ARN>}" \
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

## ENV variables required for lamba

```
REBUILD_ACTION=rebuild
REBUILD_MESSAGE=full
RIALTO_SNS_ENDPOINT=<ENDPOINT>
RIALTO_TOPIC_ARN=<ARN>
```

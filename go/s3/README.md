## LocalStack usage

This requires a LocalStack profile in `.aws/config`:
```
[profile localstack]
role_arn = arn:aws:iam::<IAM_ARN_ID>:role/devops
source_profile = default
endpoint_url = http://localhost:4566
```

1. docker daemon must be running: `orb start`
2. start LocalStack: `localstack start -d`
3. confirm it's up: `localstack status services | grep s3`
4. confirm it's working: `AWS_PROFILE=localstack aws s3 mb s3://sre-matttest`
5. copy up the test file: `AWS_PROFILE=localstack aws s3 cp data/userneeds.png s3://sre-matttest/`
6. run the go app using AWS S3: `go run . | go run . -profile test`
7. run the go app using LocalStack:`go run . -profile localstack`


## My Little Go App
This is a very basic Go application that serves a few endpoints.

## Configuration
The application can take the following environment variables to adjust behaviors of certain endpoints.

- `DYNAMODB_TABLE_NAME`: A DynamoDB Table Name
- `DYNAMODB_ITEM_HASH`: The `MyKey` value of a string Hash Key of an existing item in the DynamoDB that has a string field `MyContent`

## Endpoints

### /hello
I always like a good greeting. This endpoint should return `Hello, Moses!`

### /dynamo
When the needed environment variables are set;
Reads from a DynamoDB table (`DYNAMODB_TABLE_NAME`) by doing a Get Item with the key set to `MyKey:{"S":DYNAMODB_ITEM_HASH}` and then prints the value of a String field named `MyContent`.

## Testing

To run the test suite execute the following command:
```shell
go test ./...
```

## Local Dev

### AWS SAM Usage

1. Build the application in this directory (`GOOS=linux GOARCH=amd64 go build -o main .`).
2. Start the application using AWS SAM (`sam local start-api`).

### DynamoDB Local Setup

1. Start DynamoDB Local (`docker run --rm -p 8000:8000 amazon/dynamodb-local`).
2. Create the table inside DynamoDB Local:
```shell
AWS_REGION=us-east-1 AWS_ACCESS_KEY_ID=ALTR AWS_SECRET_ACCESS_KEY=ISGREAT \
aws dynamodb create-table \
--table-name OurTable \
--attribute-definitions AttributeName=MyKey,AttributeType=S \
--key-schema AttributeName=MyKey,KeyType=HASH \
--billing-mode PAY_PER_REQUEST \
--endpoint-url http://localhost:8000
```
3. Add an item to DynamoDB Local:
```shell
AWS_REGION=us-east-1 AWS_ACCESS_KEY_ID=ALTR AWS_SECRET_ACCESS_KEY=ISGREAT \
aws dynamodb put-item \
--table-name OurTable \
--item '{"MyKey":{"S":"this-is-a-key"}, "MyContent": {"S": "This was read from DynamoDB!"}}' \
--endpoint-url http://localhost:8000 
```
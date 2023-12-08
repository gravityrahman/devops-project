package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	EnvDynamoTableName   = "DYNAMODB_TABLE_NAME" // Env var that should contain a DynamoDB Table Name.
	EnvDynamoItemHash    = "DYNAMODB_ITEM_HASH"  // Env var that should provide the item hash key.
	DynamoItemContentKey = "MyContent"           // The field inside the DynamoDB item to return.
	EnvEnvironmentType   = "ENVIRONMENT_TYPE"    // Describe current environment type. If "DEV" uses DynamoDB Local.
	EnvDynamoDBLocalURL  = "DYNAMODB_LOCAL_URL"  // URL for DynamoDB Local (Default: http://host.docker.internal:8000)
)

func dynamoHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		dynamo, err := attemptToGetStringFromDynamo()
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
		} else {
			_, _ = fmt.Fprint(response, dynamo)
		}
	default:
		http.Error(response, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}

// attemptToGetStringFromDynamo tries to connect and retrieve a field from an item in a dynamo table.
// The jobs checks to see if the required environment variables are set and if not returns an error string.
// If the needed variables are set and we retrieve the item we return the string for the defined variable.
func attemptToGetStringFromDynamo() (string, error) {
	itemhash, okItem := os.LookupEnv(EnvDynamoItemHash)
	table, okTable := os.LookupEnv(EnvDynamoTableName)

	if !okItem || !okTable {
		return fmt.Sprintf(
			"The required variables to talk to dynamodb were not set. Please set %s & %s.",
			EnvDynamoItemHash, EnvDynamoTableName), nil
	}

	cfg, err := getAWSConfig()
	if err != nil {
		return "", err
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	item, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"MyKey": &types.AttributeValueMemberS{Value: itemhash},
		},
		TableName: aws.String(table),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get \"%s\" from \"%s\" due to: %w", itemhash, table, err)
	}

	if item.Item == nil {
		return fmt.Sprintf("No item for key %s.", itemhash), nil
	}

	value, ok := item.Item[DynamoItemContentKey]
	if !ok {
		return fmt.Sprintf("didn't find a string field(%s) in item.", DynamoItemContentKey), nil
	}

	switch v := value.(type) {
	case *types.AttributeValueMemberS:
		return v.Value, nil
	default:
		return fmt.Sprintf("Returned DynamoDB type (%T) is not String.", value), nil
	}
}

// getAWSConfig setups an AWS config with support for DynamoDB Local.
// If EnvEnvironmentType is set and equals DEV we utilize EnvDynamoDBLocalURL as a DynamoDB Local URL.
func getAWSConfig() (aws.Config, error) {
	if os.Getenv(EnvEnvironmentType) != "DEV" {
		return awsconfig.LoadDefaultConfig(context.TODO()) // nolint:wrapcheck // We're just acting as a small wrapper.
	}

	dynamoDBURL, isSet := os.LookupEnv(EnvDynamoDBLocalURL)
	if !isSet {
		dynamoDBURL = "http://host.docker.internal:8000"
	}

	dbLocal := func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && region == "us-east-1" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           dynamoDBURL,
				SigningRegion: "us-east-1",
			}, nil
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(dbLocal)

	return awsconfig.LoadDefaultConfig( // nolint:wrapcheck // We're just acting as a small wrapper.
		context.TODO(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ALTR", "ISGREAT", "")),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
}

package db

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Table struct {
	DynamoDbClient *dynamodb.Client
	TableName string
}

func SaveSubscription() {}

func GetSubscribers() {}

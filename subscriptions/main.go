package subscriptions

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	webpush "github.com/SherClockHolmes/webpush-go"
)

type SubscriptionClient struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func NewSubscriptionClient() (c *SubscriptionClient, err error) {
	tableName, err := getTableName()
	if err != nil {
		return c, err
	}
	c = &SubscriptionClient{
		TableName: tableName,
	}
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return c, err
	}
	c.DynamoDbClient = dynamodb.NewFromConfig(config)
	return c, nil
}

func getTableName() (name string, err error) {
	name = os.Getenv("SUBSCRIPTIONS_TABLE_NAME")
	if name == "" {
		err = errors.New("Environment variable unset")
	}
	return name, err
}

func (c SubscriptionClient) DoesSubscriptionExist(endpoint string) (*bool, error) {
	res, err := c.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(c.TableName),
		Key: map[string]types.AttributeValue{
			"Endpoint": &types.AttributeValueMemberS{
				Value: endpoint,
			},
		},
	})
	if (err != nil) {
		return nil, err
	}

	exists := res.Item != nil
	return &exists, nil
}

func (c SubscriptionClient) SaveSubscription(u webpush.Subscription) (err error) {
	obj, err := attributevalue.MarshalMap(u)
	if (err != nil) {
		return err
	}

	_, err = c.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: obj,
		TableName: aws.String(c.TableName),
	})
	return err
}

func (c SubscriptionClient) Unsubscribe(endpoint string) (err error) {
	_, err = c.DynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"Endpoint": &types.AttributeValueMemberS{
				Value: endpoint,
			},
		},
		TableName: aws.String(c.TableName),
	})
	return err
} 

func (c SubscriptionClient) GetSubscribers() (*[]webpush.Subscription, error) {
	metadata, err := c.DynamoDbClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(c.TableName),
	})
	if (err != nil) {
		return nil, err
	}

	slice := make([]map[string]types.AttributeValue, 0, *metadata.Table.ItemCount)
	var lastKey map[string]types.AttributeValue = nil

	for {
		res, err := c.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String(c.TableName),
			ExclusiveStartKey: lastKey,	
		})
		if (err != nil) {
			return nil, err
		}
		slice = append(slice, res.Items...)
		if (res.LastEvaluatedKey == nil) {
			break
		}
		lastKey = res.LastEvaluatedKey
	}

	var subscriptions []webpush.Subscription
	err = attributevalue.UnmarshalListOfMaps(slice, &subscriptions)
	if (err != nil) {
		return nil, err
	}
	return &subscriptions, nil
}

package db

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type db[T any] struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type DB[T any] interface {
	DeleteObject(T) error
	DoesKeyExist(map[string]types.AttributeValue) (*bool, error)
	GetObjects() (*[]T, error)
	SaveObject(T) (error)
}

func NewDB[T any]() (c *db[T], err error) {
	tableName, err := getTableName()
	if err != nil {
		return c, err
	}
	c = &db[T]{
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

// DoesObjExist returns whether an object found in the DB according to the given searchObj or an error. The searchObj is the desired object in the DB. It does not need to be complete, but should include the primary key in the database. 
func (c db[T]) DoesKeyExist(searchObj map[string]types.AttributeValue) (*bool, error) {
	res, err := c.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(c.TableName),
		Key:       searchObj,
	})
	if err != nil {
		return nil, err
	}

	exists := res.Item != nil
	return &exists, nil
}

func (c db[T]) SaveObject(obj T) (err error) {
	objAttributevalue, err := attributevalue.MarshalMap(obj)
	if err != nil {
		return err
	}

	_, err = c.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      objAttributevalue,
		TableName: aws.String(c.TableName),
	})
	return err
}

func (c db[T]) DeleteObject(obj T) (err error) {
	objAttributevalue, err := attributevalue.MarshalMap(obj)
	if err != nil {
		return err
	}

	_, err = c.DynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       objAttributevalue,
		TableName: aws.String(c.TableName),
	})
	return err
}

func (c db[T]) GetObjects() (*[]T, error) {
	metadata, err := c.DynamoDbClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(c.TableName),
	})
	if err != nil {
		return nil, err
	}

	slice := make([]map[string]types.AttributeValue, 0, *metadata.Table.ItemCount)
	var lastKey map[string]types.AttributeValue = nil

	for {
		res, err := c.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName:         aws.String(c.TableName),
			ExclusiveStartKey: lastKey,
		})
		if err != nil {
			return nil, err
		}
		slice = append(slice, res.Items...)
		if res.LastEvaluatedKey == nil {
			break
		}
		lastKey = res.LastEvaluatedKey
	}

	var objects []T
	err = attributevalue.UnmarshalListOfMaps(slice, &objects)
	if err != nil {
		return nil, err
	}
	return &objects, nil
}

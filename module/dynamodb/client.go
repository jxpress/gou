package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// Client DynamoDB用のクライアント定義
type Client interface {
	GetTableList() ([]string, error)
	CreateTable(tableName string, out interface{}) error
	DeleteTable(tableName string) error
	GetItem(tableName string, hashKey *HashKey, object interface{}) error
	GetRangeItem(tableName string, hashKey *HashKey, rangeKey *RangeKey, object interface{}) error
	GetAllItems(tableName string, hashKey *HashKey, objects interface{}) error
	GetPaginationFirstItems(tableName string, hashKey *HashKey, limit int64, objects interface{}) error
	ScanAll(tableName string, objects interface{}) error
	Put(tableName string, object interface{}) error
	Update(tableName string, hashKey string, object interface{}) error
	DeleteItem(tableName string, hashKey *HashKey, rangeKey *RangeKey) error
}

type dynamoDBClient struct {
	db *dynamo.DB
	Client
}

// HashKey DynamoDBクライアントからGetItemする際に指定するhash key
type HashKey struct {
	Name  string
	Value interface{}
}

// RangeKey DynamoDBクライアントからGetItemする際に指定するrange key
type RangeKey struct {
	Name  string
	Value interface{}
}

// NewClient DynamoDB用のクライアントを生成する
func NewClient(config *aws.Config) (Client, error) {
	session, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	db := dynamo.New(session, config)
	return &dynamoDBClient{
		db: db,
	}, nil
}

func (c *dynamoDBClient) CreateTable(tableName string, out interface{}) error {
	return c.db.CreateTable(tableName, out).Run()
}

func (c *dynamoDBClient) DeleteTable(tableName string) error {
	return c.db.Table(tableName).DeleteTable().Run()
}

func (c *dynamoDBClient) GetTableList() ([]string, error) {
	return c.db.ListTables().All()
}

func (c *dynamoDBClient) GetItem(tableName string, hashKey *HashKey, object interface{}) error {
	return c.getTable(tableName).Get(hashKey.Name, hashKey.Value).One(object)
}

func (c *dynamoDBClient) GetRangeItem(tableName string, hashKey *HashKey, rangeKey *RangeKey, object interface{}) error {
	return c.getTable(tableName).Get(hashKey.Name, hashKey.Value).Range(rangeKey.Name, dynamo.Equal, rangeKey.Value).One(object)
}

func (c *dynamoDBClient) GetAllItems(tableName string, hashKey *HashKey, objects interface{}) error {
	return c.getTable(tableName).Get(hashKey.Name, hashKey.Value).All(objects)
}

func (c *dynamoDBClient) GetPaginationFirstItems(tableName string, hashKey *HashKey, limit int64, objects interface{}) error {
	table := c.getTable(tableName).Get(hashKey.Name, hashKey.Value)
	_, err := table.Order(false).Limit(limit).AllWithLastEvaluatedKey(objects)
	return err
}

func (c *dynamoDBClient) ScanAll(tableName string, objects interface{}) error {
	return c.getTable(tableName).Scan().All(objects)
}

func (c *dynamoDBClient) Put(tableName string, object interface{}) error {
	return c.getTable(tableName).Put(object).Run()
}

func (c *dynamoDBClient) Update(tableName string, hashKey string, object interface{}) error {
	return c.getTable(tableName).Update(hashKey, object).Run()
}

func (c *dynamoDBClient) DeleteItem(tableName string, hashKey *HashKey, rangeKey *RangeKey) error {
	return c.getTable(tableName).Delete(hashKey.Name, hashKey.Value).Range(rangeKey.Name, rangeKey.Value).Run()
}

func (c *dynamoDBClient) getTable(tableName string) dynamo.Table {
	return c.db.Table(tableName)
}

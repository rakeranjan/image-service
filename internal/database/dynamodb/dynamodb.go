package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rakeranjan/image-service/internal/configuration"
	"github.com/rakeranjan/image-service/utils"
)

type DataBase struct {
	db            *dynamodb.Client
	configuration configuration.Configuration
}

var (
	instance *DataBase
	once     sync.Once
)

func NewDataBase() *DataBase {
	configObj, _ := configuration.NewConfiguration()
	var err error
	once.Do(func() {
		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(configObj.Region), config.WithBaseEndpoint(configObj.AwsBaseEndPoint))
		if err != nil {
			fmt.Println("error while loading config", err)
		}
		svc := dynamodb.NewFromConfig(config)
		instance = &DataBase{
			db:            svc,
			configuration: configObj,
		}
	})
	if err != nil {
		return nil
	}

	return instance
}

func (d *DataBase) GetDB() *dynamodb.Client {
	return d.db
}

func (d *DataBase) CreateUserTable() {
	ok, err := d.TableExists(context.TODO(), utils.USERS_TABLE)
	if ok {
		log.Println("user table is already present")
		return
	}
	table, err := d.db.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("phoneNumber"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("userName"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("phoneNumber"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String(utils.USERS_TABLE),
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("userName-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("userName"),
						KeyType:       types.KeyTypeHash, // Partition key for GSI
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll, // Include all attributes in the index
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1000),
					WriteCapacityUnits: aws.Int64(100),
				},
			},
			{
				IndexName: aws.String("phoneNumber-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("phoneNumber"),
						KeyType:       types.KeyTypeHash, // Partition key for GSI
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1000),
					WriteCapacityUnits: aws.Int64(100),
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1000),
			WriteCapacityUnits: aws.Int64(100),
		},
	})

	if err != nil {
		fmt.Println("error while creating table", err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(d.db)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(utils.USERS_TABLE)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc := table.TableDescription
		fmt.Println(tableDesc, table)
	}
}

func (d *DataBase) CreateImageTable() {
	ok, err := d.TableExists(context.TODO(), utils.IMAGE_TABLE)
	if ok {
		return
	}
	table, err := d.db.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("imageId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("userPhoneNumber"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("userId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("userId"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String(utils.IMAGE_TABLE),
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("imageId-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("imageId"),
						KeyType:       types.KeyTypeHash, // Partition key for GSI
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll, // Include all attributes in the index
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1000),
					WriteCapacityUnits: aws.Int64(100),
				},
			},
			{
				IndexName: aws.String("userPhoneNumber-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("userPhoneNumber"),
						KeyType:       types.KeyTypeHash, // Partition key for GSI
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1000),
					WriteCapacityUnits: aws.Int64(100),
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1000),
			WriteCapacityUnits: aws.Int64(100),
		},
	})

	if err != nil {
		fmt.Println("error while creating table", err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(d.db)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(utils.IMAGE_TABLE)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc := table.TableDescription
		fmt.Println(tableDesc, table)
	}
}

func (d *DataBase) TableExists(ctx context.Context, tableName string) (bool, error) {
	exists := true
	_, err := d.db.DescribeTable(
		ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist.\n", tableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", tableName, err)
		}
		exists = false
	}
	return exists, err
}

func (d *DataBase) DeleteTable(ctx context.Context, tableName string) error {
	_, err := d.db.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName)})
	if err != nil {
		log.Printf("Couldn't delete table %v. Here's why: %v\n", tableName, err)
	}
	return err
}

package objectdatabase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/rakeranjan/image-service/internal/configuration"
)

func init() {
	configObj, _ = configuration.NewConfiguration()
}

var configObj *configuration.Config

type ObjectStorage struct {
	client        *s3.Client
	configuration *configuration.Config
}

var (
	instance *ObjectStorage

	once sync.Once
)

func NewObjectStorage() *ObjectStorage {

	var err error
	once.Do(func() {
		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(configObj.Region), config.WithBaseEndpoint(configObj.AwsBaseEndPoint))
		if err != nil {
			fmt.Println("error while loading config", err)
		}
		config.EndpointResolver = CustomEndpointResolver{}
		client := s3.NewFromConfig(config)
		instance = &ObjectStorage{
			client:        client,
			configuration: configObj,
		}
	})
	if err != nil {
		return nil
	}
	return instance
}

func (o *ObjectStorage) GetClient() *s3.Client {
	return o.client
}

func (o *ObjectStorage) CreateBucket(ctx context.Context, name string) error {
	_, err := o.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(o.configuration.Region),
		},
	})

	if err != nil {
		fmt.Println("error: ", err)
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			log.Printf("You already own bucket %s.\n", name)
			err = owned
		} else if errors.As(err, &exists) {
			log.Printf("Bucket %s already exists.\n", name)
			err = exists
		}
	} else {
		err = s3.NewBucketExistsWaiter(o.client).Wait(
			ctx, &s3.HeadBucketInput{Bucket: aws.String(name)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for bucket %s to exist.\n", name)
		}
	}
	return err
}

func (o *ObjectStorage) ListBuckets(ctx context.Context) ([]string, error) {
	result, err := o.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "AccessDenied" {
			fmt.Println("You don't have permission to list buckets for this account.")
			return nil, fmt.Errorf("You don't have permission to list buckets for this account.\n")
		} else {
			return nil, fmt.Errorf("Couldn't list buckets for your account. Here's why: %v\n", err)
		}
	}
	buckets := make([]string, 0)
	for _, v := range result.Buckets {
		buckets = append(buckets, *v.Name)
	}
	return buckets, nil
}

type CustomEndpointResolver struct {
}

func (c CustomEndpointResolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) {
	configObj, _ := configuration.NewConfiguration()
	if service == s3.ServiceID && region == configObj.Region {
		return aws.Endpoint{
			URL:               configObj.AwsBaseEndPoint,
			SigningRegion:     configObj.Region,
			HostnameImmutable: true, // Forces path-style addressing by keeping hostname unchanged
		}, nil
	}
	return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
}

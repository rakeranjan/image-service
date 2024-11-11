package configuration

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

const (
	SECRET_VALUE         = "SECRET_VALUE"
	MIN_IMAGE_SIZE_IN_MB = "MIN_IMAGE_SIZE_IN_MB"
	MAX_IMAGE_SIZE_IN_MB = "MAX_IMAGE_SIZE_IN_MB"
	REGION               = "REGION"
	AWS_BASE_END_POINT   = "BASE_END_POINT"
	PROCESSING_QUEUE_URL = "PROCESSING_QUEUE_URL"
	PROCESSED_QUEUE_URL  = "PROCESSED_QUEUE_URL"
	MODULE               = "MODULE"
	PORT                 = "PORT"
)

type Config struct {
	SecretValue        string
	MinImageSizeInMB   int
	MaxImageSizeInMB   int
	Region             string
	AwsBaseEndPoint    string
	ProcessingQueueURL string
	ProcessedQueueURL  string
	Module             string
	Port               string
}

var (
	instance *Config
	once     sync.Once
)

func NewConfiguration() (*Config, error) {
	var err error
	once.Do(func() {
		secretValue := os.Getenv(SECRET_VALUE)
		if secretValue == "" {
			err = fmt.Errorf("%s is empty", SECRET_VALUE)
		}
		module := os.Getenv(MODULE)
		if module == "" {
			err = fmt.Errorf("%s is empty", MODULE)
		}
		port := os.Getenv(PORT)
		if port == "" {
			err = fmt.Errorf("%s is empty", PORT)
		}
		// Port               string

		processingQueueURL := os.Getenv(PROCESSING_QUEUE_URL)
		if processingQueueURL == "" {
			err = fmt.Errorf("%s is empty", PROCESSING_QUEUE_URL)
		}

		processedQueueURL := os.Getenv(PROCESSED_QUEUE_URL)
		if processedQueueURL == "" {
			err = fmt.Errorf("%s is empty", PROCESSED_QUEUE_URL)
		}

		region := os.Getenv(REGION)
		if region == "" {
			err = fmt.Errorf("%s is empty", REGION)
		}

		awsBaseEndPoint := os.Getenv(AWS_BASE_END_POINT)
		if awsBaseEndPoint == "" {
			err = fmt.Errorf("%s is empty", AWS_BASE_END_POINT)
		}

		minImageSizeInMBData := os.Getenv(MIN_IMAGE_SIZE_IN_MB)
		if minImageSizeInMBData == "" {
			err = fmt.Errorf("%s is empty", MIN_IMAGE_SIZE_IN_MB)
		}

		minImageSizeInMB, err := strconv.Atoi(minImageSizeInMBData)
		if err != nil {
			err = fmt.Errorf("%s is not a valid value: %s", MIN_IMAGE_SIZE_IN_MB, minImageSizeInMBData)
		}

		maxImageSizeInMBData := os.Getenv(MAX_IMAGE_SIZE_IN_MB)
		if maxImageSizeInMBData == "" {
			err = fmt.Errorf("%s is empty", MAX_IMAGE_SIZE_IN_MB)
		}

		maxImageSizeInMB, err := strconv.Atoi(maxImageSizeInMBData)
		if err != nil {
			err = fmt.Errorf("%s is not a valid value: %s", MAX_IMAGE_SIZE_IN_MB, maxImageSizeInMBData)
		}
		instance = &Config{
			SecretValue:        secretValue,
			MinImageSizeInMB:   minImageSizeInMB,
			MaxImageSizeInMB:   maxImageSizeInMB,
			Region:             region,
			AwsBaseEndPoint:    awsBaseEndPoint,
			ProcessingQueueURL: processingQueueURL,
			ProcessedQueueURL:  processedQueueURL,
			Module:             module,
			Port:               port,
		}
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *Config) GetSecretValue() string {
	return c.SecretValue
}

func (c *Config) GetMinImageSizeInMB() int {
	return c.MinImageSizeInMB
}

func (c *Config) GetMaxImageSizeInMB() int {
	return c.MaxImageSizeInMB
}

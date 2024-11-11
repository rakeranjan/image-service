package models

import "time"

type ImageMetaData struct {
	FileName        string    `dynamodbav:"fileName"`
	ImageId         string    `dynamodbav:"imageId"`
	CreatedAt       time.Time `dynamodbav:"createdAt"`
	UserId          string    `dynamodbav:"userId"`
	UserPhoneNumber string    `dynamodbav:"userPhoneNumber"`
	SizeInKb        int       `dynamodbav:"sizeInKb"`
	IsProcessed     bool      `dynamodbav:"isProcessed"`
}

func (i *ImageMetaData) GetObjectKey() string {
	return i.ImageId
}

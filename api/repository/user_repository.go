package repository

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/rakeranjan/image-service/api/models"
	db "github.com/rakeranjan/image-service/internal/database/dynamodb"
	"github.com/rakeranjan/image-service/utils"
)

type UserRepository struct {
	db dynamodb.Client
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: *db.NewDataBase().GetDB(),
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	userPresent, err := u.GetUserByPhoneNumber(ctx, user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if userPresent != nil {
		return userPresent, nil
	}
	user.ID = uuid.NewString()
	items, err := attributevalue.MarshalMap(user)
	if err != nil {
		log.Println("error while marshling user, err:", err)
		return nil, err
	}
	_, err = u.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(utils.USERS_TABLE),
		Item:      items,
	})
	if err != nil {
		log.Println("error while saving user, err:", err)
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserByPhoneNumber(ctx context.Context, phoneNmber string) (*models.User, error) {
	opt, err := u.db.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"phoneNumber": &types.AttributeValueMemberS{
				Value: *aws.String(phoneNmber),
			},
		},
		TableName: aws.String(utils.USERS_TABLE),
	})
	var user models.User
	if err != nil {
		log.Println("error while searching user, err:", err)
		return nil, err
	}
	if opt.Item == nil {
		return nil, err
	}
	attributevalue.UnmarshalMap(opt.Item, &user)

	return &user, nil
}

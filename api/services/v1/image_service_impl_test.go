package v1

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/rakeranjan/image-service/api/models"
	"github.com/stretchr/testify/mock"
)

func TestImageServiceImpl_Upload(t *testing.T) {
	type mockDetails struct {
		mockObj                      *MockImageRepository
		UploadToProcessingArg        *models.ImageMetaData
		SendToSqsForProcessingArg    *models.ImageMetaData
		SaveImageMetaDataErr         error
		UploadToProcessingArgErr     error
		SendToSqsForProcessingArgErr error
	}
	type args struct {
		ctx        context.Context
		user       *models.User
		fileHeader *multipart.FileHeader
	}
	tests := []struct {
		name        string
		mockDetails *mockDetails
		args        args
		want        *models.ImageMetaData
		wantErr     bool
	}{
		{
			name: "negative with nil file",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:          "123",
					FirstName:   "as",
					LastName:    "bv",
					UserName:    "ac.df",
					PhoneNumber: "9876543210",
				},
				fileHeader: nil,
			},
			want: nil,
			mockDetails: &mockDetails{
				mockObj: new(MockImageRepository),
			},
			wantErr: true,
		},
		{
			name: "negative with nil user",
			args: args{
				ctx:        context.Background(),
				user:       nil,
				fileHeader: &multipart.FileHeader{Size: 12},
			},
			want: nil,
			mockDetails: &mockDetails{
				mockObj: new(MockImageRepository),
			},
			wantErr: true,
		},
		{
			name: "negative with saveImageMetaData error",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:          "123",
					FirstName:   "as",
					LastName:    "bv",
					UserName:    "ac.df",
					PhoneNumber: "9876543210",
				},
				fileHeader: &multipart.FileHeader{Size: 12},
			},
			want: nil,
			mockDetails: &mockDetails{
				mockObj:              new(MockImageRepository),
				SaveImageMetaDataErr: errors.New("some error"),
			},
			wantErr: true,
		},
		{
			name: "negative with UploadToProcessing error",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:          "123",
					FirstName:   "as",
					LastName:    "bv",
					UserName:    "ac.df",
					PhoneNumber: "9876543210",
				},
				fileHeader: &multipart.FileHeader{Size: 12},
			},
			want: nil,
			mockDetails: &mockDetails{
				mockObj:                  new(MockImageRepository),
				SaveImageMetaDataErr:     nil,
				UploadToProcessingArgErr: errors.New("some error"),
			},
			wantErr: true,
		},
		{
			name: "negative with SendToSqsForProcessing error",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:          "123",
					FirstName:   "as",
					LastName:    "bv",
					UserName:    "ac.df",
					PhoneNumber: "9876543210",
				},
				fileHeader: &multipart.FileHeader{Size: 12},
			},
			want: nil,
			mockDetails: &mockDetails{
				mockObj:                      new(MockImageRepository),
				SaveImageMetaDataErr:         nil,
				UploadToProcessingArgErr:     nil,
				SendToSqsForProcessingArgErr: errors.New("Some err"),
			},
			wantErr: true,
		},
		{
			name: "positive with no error",
			args: args{
				ctx: context.Background(),
				user: &models.User{
					ID:          "123",
					FirstName:   "as",
					LastName:    "bv",
					UserName:    "ac.df",
					PhoneNumber: "9876543210",
				},
				fileHeader: &multipart.FileHeader{Size: 12},
			},
			want: &models.ImageMetaData{
				UserId:          "123",
				UserPhoneNumber: "9876543210",
			},
			mockDetails: &mockDetails{
				mockObj:                      new(MockImageRepository),
				SaveImageMetaDataErr:         nil,
				UploadToProcessingArgErr:     nil,
				SendToSqsForProcessingArgErr: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := ImageServiceImpl{
				imageRepository: tt.mockDetails.mockObj,
			}
			// SaveImageMetaData(ctx context.Context, data *models.ImageMetaData) error
			tt.mockDetails.mockObj.On("SaveImageMetaData", tt.args.ctx, mock.Anything).Return(tt.mockDetails.SaveImageMetaDataErr)
			// UploadToProcessing(ctx context.Context, metaData *models.ImageMetaData, 	fileHeader *multipart.FileHeader) error
			tt.mockDetails.mockObj.On("UploadToProcessing", tt.args.ctx, mock.Anything, tt.args.fileHeader).Return(tt.mockDetails.UploadToProcessingArgErr)
			// SendToSqsForProcessing(ctx context.Context, imageMetaData *models.ImageMetaData) error
			tt.mockDetails.mockObj.On("SendToSqsForProcessing", tt.args.ctx, mock.Anything).Return(tt.mockDetails.SendToSqsForProcessingArgErr)
			got, err := i.Upload(tt.args.ctx, tt.args.user, tt.args.fileHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageServiceImpl.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got.UserId != tt.want.UserId || got.UserPhoneNumber != tt.want.UserPhoneNumber) {
				t.Errorf("ImageServiceImpl.Upload() = %v, want %v", got, tt.want)
			}
		})
	}
}

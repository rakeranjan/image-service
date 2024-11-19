package v1

import (
	"context"
	"errors"
	"mime/multipart"
	"reflect"
	"testing"
	"time"

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

func TestImageServiceImpl_GetByID(t *testing.T) {
	type mockDetails struct {
		mockObj           *MockImageRepository
		mockImageResponse *models.ImageResponse
		mockErr           error
	}
	type args struct {
		ctx     context.Context
		user    *models.User
		imageID string
	}
	tests := []struct {
		name        string
		args        args
		mockDetails mockDetails
		want        *models.ImageResponse
		wantErr     bool
	}{
		{
			name: "positive",
			mockDetails: mockDetails{
				mockObj: new(MockImageRepository),
				mockImageResponse: &models.ImageResponse{
					ImageMetaData: models.ImageMetaData{
						FileName:        "a.jpg",
						ImageId:         "123-abc",
						CreatedAt:       time.Time{},
						UserId:          "123-abc",
						UserPhoneNumber: "9876543210",
						SizeInKb:        29,
						IsProcessed:     false,
						Height:          12,
						Width:           34,
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				user: &models.User{
					ID:          "123-1bc",
					PhoneNumber: "9876543210",
				},
				imageID: "found",
			},
			want: &models.ImageResponse{
				ImageMetaData: models.ImageMetaData{
					FileName:        "a.jpg",
					ImageId:         "123-abc",
					CreatedAt:       time.Time{},
					UserId:          "123-abc",
					UserPhoneNumber: "9876543210",
					SizeInKb:        29,
					IsProcessed:     false,
					Height:          12,
					Width:           34,
				},
			},
			wantErr: false,
		},
		{
			name: "with error",
			mockDetails: mockDetails{
				mockObj:           new(MockImageRepository),
				mockImageResponse: nil,
				mockErr:           errors.New("some error"),
			},
			args: args{
				ctx: context.TODO(),
				user: &models.User{
					ID:          "123-1bc",
					PhoneNumber: "9876543210",
				},
				imageID: "not-found",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := ImageServiceImpl{
				imageRepository: tt.mockDetails.mockObj,
			}
			tt.mockDetails.mockObj.On("GetImageMetaDataByImageID", tt.args.ctx, tt.args.user, tt.args.imageID).Return(tt.mockDetails.mockImageResponse, tt.mockDetails.mockErr)
			got, err := i.GetByID(tt.args.ctx, tt.args.user, tt.args.imageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageServiceImpl.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImageServiceImpl.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageServiceImpl_List(t *testing.T) {
	type mockDetails struct {
		mockObj      *MockImageRepository
		mockResponse string
		mockErr      error
	}
	type args struct {
		ctx  context.Context
		user *models.User
	}
	tests := []struct {
		name        string
		mockDetails mockDetails
		args        args
		want        string
		wantErr     bool
	}{
		{
			name: "positive",
			mockDetails: mockDetails{
				mockObj:      new(MockImageRepository),
				mockResponse: "./some-file.jpg",
				mockErr:      nil,
			},
			args: args{
				ctx: context.TODO(),
				user: &models.User{
					ID:          "123-1bc",
					PhoneNumber: "9876543210",
				},
			},
			want:    "./some-file.jpg",
			wantErr: false,
		},
		{
			name: "with error",
			mockDetails: mockDetails{
				mockObj:      new(MockImageRepository),
				mockResponse: "",
				mockErr:      errors.New("some error"),
			},
			args: args{
				ctx: context.TODO(),
				user: &models.User{
					ID:          "123-abc",
					PhoneNumber: "9876543210",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := ImageServiceImpl{
				imageRepository: tt.mockDetails.mockObj,
			}
			tt.mockDetails.mockObj.On("GetAllImageMetaData", tt.args.ctx, tt.args.user).Return(tt.mockDetails.mockResponse, tt.mockDetails.mockErr)
			got, err := i.List(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageServiceImpl.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImageServiceImpl.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageServiceImpl_Update(t *testing.T) {
	type mockDetails struct {
		mockObj                           *MockImageRepository
		UploadToProcessingResponse        *models.ImageMetaData
		SendToSqsForProcessingResponse    *models.ImageMetaData
		GetImageMetaDataByImageIDResponse *models.ImageResponse
		DeleteProcessedObjextErr          error
		SaveImageMetaDataErr              error
		UploadToProcessingArgErr          error
		SendToSqsForProcessingArgErr      error
		GetImageMetaDataByImageIDErr      error
	}
	type args struct {
		ctx        context.Context
		user       *models.User
		imageID    string
		fileHeader *multipart.FileHeader
	}
	tests := []struct {
		name        string
		mockDetails mockDetails
		args        args
		want        *models.ImageMetaData
		wantErr     bool
	}{
		{
			name: "positive",
			mockDetails: mockDetails{
				mockObj: new(MockImageRepository),
				UploadToProcessingResponse: &models.ImageMetaData{
					FileName:        "some-file.jpg",
					ImageId:         "123-abc",
					CreatedAt:       time.Time{},
					UserId:          "123-abc",
					UserPhoneNumber: "9876543210",
					SizeInKb:        10,
					IsProcessed:     false,
					Height:          10,
					Width:           10,
				},
				GetImageMetaDataByImageIDResponse: &models.ImageResponse{
					ImageMetaData: models.ImageMetaData{
						FileName:        "some-file.jpg",
						ImageId:         "123-abc",
						CreatedAt:       time.Time{},
						UserId:          "123-abc",
						UserPhoneNumber: "9876543210",
						SizeInKb:        10,
						IsProcessed:     false,
						Height:          10,
						Width:           10,
					},
				},
				SendToSqsForProcessingResponse: &models.ImageMetaData{
					FileName:        "some-file.jpg",
					ImageId:         "123-abc",
					CreatedAt:       time.Time{},
					UserId:          "123-abc",
					UserPhoneNumber: "9876543210",
					SizeInKb:        10,
					IsProcessed:     false,
					Height:          10,
					Width:           10,
				},
				SaveImageMetaDataErr:         nil,
				UploadToProcessingArgErr:     nil,
				SendToSqsForProcessingArgErr: nil,
				GetImageMetaDataByImageIDErr: nil,
			},
			args: args{
				ctx: context.TODO(),
				user: &models.User{
					ID:          "123-abc",
					PhoneNumber: "9876543210",
				},
				imageID:    "found",
				fileHeader: &multipart.FileHeader{},
			},
			want: &models.ImageMetaData{
				FileName:        "some-file.jpg",
				ImageId:         "123-abc",
				CreatedAt:       time.Time{},
				UserId:          "123-abc",
				UserPhoneNumber: "9876543210",
				SizeInKb:        10,
				IsProcessed:     false,
				Height:          10,
				Width:           10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := ImageServiceImpl{
				imageRepository: tt.mockDetails.mockObj,
			}
			tt.mockDetails.mockObj.On("DeleteProcessedObjext", tt.args.ctx, mock.Anything).Return(tt.mockDetails.DeleteProcessedObjextErr)
			tt.mockDetails.mockObj.On("GetImageMetaDataByImageID", tt.args.ctx, tt.args.user, tt.args.imageID).Return(tt.mockDetails.GetImageMetaDataByImageIDResponse, tt.mockDetails.GetImageMetaDataByImageIDErr)
			tt.mockDetails.mockObj.On("SaveImageMetaData", tt.args.ctx, mock.Anything).Return(tt.mockDetails.SaveImageMetaDataErr)
			tt.mockDetails.mockObj.On("UploadToProcessing", tt.args.ctx, mock.Anything, tt.args.fileHeader).Return(tt.mockDetails.UploadToProcessingArgErr)
			tt.mockDetails.mockObj.On("SendToSqsForProcessing", tt.args.ctx, mock.Anything).Return(tt.mockDetails.SendToSqsForProcessingArgErr)
			got, err := i.Update(tt.args.ctx, tt.args.user, tt.args.imageID, tt.args.fileHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageServiceImpl.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImageServiceImpl.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageServiceImpl_Delete(t *testing.T) {
	type fields struct {
		imageRepository ImageRepository
	}
	type args struct {
		ctx     context.Context
		user    *models.User
		imageID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := ImageServiceImpl{
				imageRepository: tt.fields.imageRepository,
			}

			if got := i.Delete(tt.args.ctx, tt.args.user, tt.args.imageID); got != tt.want {
				t.Errorf("ImageServiceImpl.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

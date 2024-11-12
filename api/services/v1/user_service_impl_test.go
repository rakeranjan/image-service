package v1

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/rakeranjan/image-service/api/models"
	"github.com/rakeranjan/image-service/internal/configuration"
)

func TestUserServiceImplV1_Create(t *testing.T) {
	type mockData struct {
		mockUserRepository *MockUserReposity
		mockExpectedUser   *models.User
		mockError          error
	}
	type fields struct {
		userRepository UserRepository
		conf           *configuration.Config
	}
	type args struct {
		ctx  context.Context
		user *models.User
	}
	tests := []struct {
		name     string
		fields   fields
		mockData mockData
		args     args
		want     *models.UserResponse
		wantErr  bool
	}{
		{
			name: "positive",
			fields: fields{
				conf: &configuration.Config{
					SecretValue: "MY_VALUE",
				},
			},
			args: args{
				ctx: context.TODO(),
				user: &models.User{
					FirstName:   "a",
					LastName:    "b",
					UserName:    "an.cn",
					Password:    "sdk.sd",
					PhoneNumber: "9876543210",
				},
			},
			want: &models.UserResponse{
				ID:          "123",
				FirstName:   "a",
				LastName:    "b",
				UserName:    "an.cn",
				PhoneNumber: "9876543210",
			},
			mockData: mockData{
				mockUserRepository: new(MockUserReposity),
				mockExpectedUser: &models.User{
					ID:          "123",
					FirstName:   "a",
					LastName:    "b",
					UserName:    "an.cn",
					Password:    "sdk.sd",
					PhoneNumber: "9876543210",
				},
				mockError: nil,
			},
			wantErr: false,
		},
		{
			name: "Negative with an error",
			fields: fields{
				conf: &configuration.Config{
					SecretValue: "MY_VALUE",
				},
			},
			args: args{
				ctx: context.TODO(),
				user: &models.User{
					LastName:    "b",
					UserName:    "an.cn",
					Password:    "sdk.sd",
					PhoneNumber: "9876543210",
				},
			},
			want: nil,
			mockData: mockData{
				mockUserRepository: new(MockUserReposity),
				mockExpectedUser:   nil,
				mockError:          errors.New("Some error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockData.mockUserRepository.On("CreateUser", tt.args.ctx, tt.args.user).Return(tt.mockData.mockExpectedUser, tt.mockData.mockError)
			i := UserServiceImplV1{
				userRepository: tt.mockData.mockUserRepository,
				conf:           tt.fields.conf,
			}
			if tt.mockData.mockExpectedUser != nil {
				tt.want.AccessToken, _ = tt.mockData.mockExpectedUser.Encrypt(tt.fields.conf.SecretValue)
			}
			got, err := i.Create(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserServiceImplV1.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserServiceImplV1.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

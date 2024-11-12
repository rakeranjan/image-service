package v1

import (
	"context"
	"reflect"
	"testing"

	"github.com/rakeranjan/image-service/api/models"
	"github.com/rakeranjan/image-service/internal/configuration"
)

func TestUserServiceImplV1_Create(t *testing.T) {
	type mockDetails struct {
		mockExpectedUser *models.UserResponse
		mockError        error
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
		name    string
		fields  fields
		args    args
		want    *models.UserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := UserServiceImplV1{
				userRepository: tt.fields.userRepository,
				conf:           tt.fields.conf,
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

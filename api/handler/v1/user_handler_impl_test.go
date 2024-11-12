package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rakeranjan/image-service/api/middleware/validators"
	"github.com/rakeranjan/image-service/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandlerImpl_Create(t *testing.T) {
	type mockDetails struct {
		mockObject *MockUserService
		// *models.UserResponse, error
		mockUserResponse *models.UserResponse
		mockError        error
	}
	type fields struct {
		UserService UserServiceV1
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mockDetails mockDetails
		wantStatus  int
		user        *models.User
	}{
		{
			name: "negative with invalid userName",
			mockDetails: mockDetails{
				mockObject: new(MockUserService),
			},
			args: args{
				c: &gin.Context{},
			},
			user: &models.User{
				FirstName:   "a",
				LastName:    "b",
				UserName:    "abnd.@",
				Password:    "asdck",
				PhoneNumber: "9876543210",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative with no userName",
			mockDetails: mockDetails{
				mockObject: new(MockUserService),
			},
			args: args{
				c: &gin.Context{},
			},
			user: &models.User{
				FirstName:   "a",
				LastName:    "b",
				Password:    "asdck",
				PhoneNumber: "9876543210",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative with no phoneNumber",
			mockDetails: mockDetails{
				mockObject: new(MockUserService),
			},
			args: args{
				c: &gin.Context{},
			},
			user: &models.User{
				FirstName: "a",
				LastName:  "b",
				UserName:  "abnd.q",
				Password:  "asdck",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "positive",
			mockDetails: mockDetails{
				mockObject: new(MockUserService),
			},
			args: args{
				c: &gin.Context{},
			},
			user: &models.User{
				FirstName:   "a",
				LastName:    "b",
				UserName:    "abnd.a",
				Password:    "asdck",
				PhoneNumber: "9876543210",
			},
			wantStatus: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandlerImpl{
				UserService: tt.mockDetails.mockObject,
			}
			g := setupUSerRoutes(u)
			w := httptest.NewRecorder()
			userJson, _ := json.Marshal(tt.user)
			req, _ := http.NewRequest("POST", "/v1/user", strings.NewReader(string(userJson)))
			tt.mockDetails.mockObject.On("Create", mock.Anything, tt.user).Return(tt.mockDetails.mockUserResponse, tt.mockDetails.mockError)
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestUserHandler_Create_BindJSONError(t *testing.T) {
	mockUserService := new(MockUserService)
	userHandler := &UserHandlerImpl{
		UserService: mockUserService,
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("userNameFormat", validators.ValidateUserName)
		v.RegisterValidation("userNameFormat", validators.ValidateUserName)
		v.RegisterValidation("phoneNumberFormat", validators.ValidatePhoneNumber)
	}
	router.POST("/users", userHandler.Create)
	invalidUserName := `{
    "firstName": "a",
    "lastName": "b",
    "userName": "cat.ran@",
    "password": "cat.ran@",
    "phoneNumber": "1234567890"
}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(invalidUserName)))
	req.Header.Set("Content-Type", "application/json")
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)
	assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
}

func setupUSerRoutes(userHandlerImpl *UserHandlerImpl) *gin.Engine {
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("userNameFormat", validators.ValidateUserName)
		v.RegisterValidation("phoneNumberFormat", validators.ValidatePhoneNumber)
	}
	v1 := r.Group("v1")
	{
		v1.POST("/user", userHandlerImpl.Create)
	}
	return r
}

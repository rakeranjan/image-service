package models

type UserResponse struct {
	User
	AccessToken string `json:"access_token"`
}

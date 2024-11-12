package models

type UserResponse struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	UserName    string `json:"userName"`
	PhoneNumber string `json:"phoneNumber"`
	AccessToken string `json:"access_token"`
}

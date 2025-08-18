package model

import "time"


type UserData struct {
	ID  int  `json:"id"`
	FristName  string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email    string  `json:"email"`
	Phone   *string  `json:"phone"`
	AvatarUrl  *string `json:"avatarUrl"`
	PhoneVerifiedAt  *time.Time `json:"PhoneVerified"`
	Status string  `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}



type GoogleCallback struct {
	Email  string `json:"email"`
    FamilyName string `json:"family_name"`
    GiveName  string `json:"given_name"`
    Name   string `json:"name"`
    Picture string `json:"picture"`
	VerifiedEmail string `json:"verified_email"`
}
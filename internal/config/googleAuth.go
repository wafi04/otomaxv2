package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


var GoogleOauthConfig = &oauth2.Config{
	ClientID:     "",
	ClientSecret: "",
	RedirectURL:  "http://localhost:8080/api/auth/google/callback", 
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}
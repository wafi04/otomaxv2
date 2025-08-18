package repository

import (
	"database/sql"
)


type AuthRepository struct {
	repo *sql.DB
}


func NewAuthRepository(repo *sql.DB) *AuthRepository{
	return &AuthRepository{
		repo: repo,
	}
}


// func (repo *AuthRepository) Create(){
// 	state := "randomstate123" 
// 	url := config.GoogleOauthConfig.AuthCodeURL(state)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }
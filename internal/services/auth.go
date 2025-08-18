package services

import "github.com/wafi04/otomaxv2/internal/repository"

type AuthService struct {
	repo *repository.AuthRepository
}


func NewAuthService(repo *repository.AuthRepository) *AuthService{
	return &AuthService{
		repo: repo,
	}
}


func (s *AuthService) GetCallback(){
	
}
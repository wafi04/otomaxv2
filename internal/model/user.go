package model

type UserRole string

const (
	RolePlatinum UserRole = "PLATINUM"
	RoleAdmin    UserRole = "ADMIN"
	RoleMember   UserRole = "MEMBER"
)

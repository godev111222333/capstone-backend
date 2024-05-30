package model

import "time"

type RoleID int

const (
	RoleIDAdmin    RoleID = 1
	RoleIDCustomer RoleID = 2
	RoleIDPartner  RoleID = 3
)

type Role struct {
	ID        int
	RoleName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

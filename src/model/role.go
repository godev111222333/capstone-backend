package model

import "time"

type RoleID int

const (
	RoleIDAdmin    RoleID = 1
	RoleIDCustomer RoleID = 2
	RoleIDPartner  RoleID = 3

	RoleCodeAdmin    = "AD"
	RoleCodeCustomer = "CS"
	RoleCodePartner  = "PN"
)

type Role struct {
	ID        int
	RoleName  string
	RoleCode  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

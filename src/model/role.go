package model

import "time"

type RoleID int

const (
	RoleIDAdmin      RoleID = 1
	RoleIDCustomer   RoleID = 2
	RoleIDPartner    RoleID = 3
	RoleIDTechnician RoleID = 4

	RoleCodeAdmin      = "AD"
	RoleCodeCustomer   = "CS"
	RoleCodePartner    = "PN"
	RoleCodeTechnician = "TN"
	RoleNameAdmin      = "admin"
	RoleNamePartner    = "partner"
	RoleNameCustomer   = "customer"
	RoleNameTechnician = "technician"
)

type Role struct {
	ID        int       `json:"id"`
	RoleName  string    `json:"role_name"`
	RoleCode  string    `json:"role_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

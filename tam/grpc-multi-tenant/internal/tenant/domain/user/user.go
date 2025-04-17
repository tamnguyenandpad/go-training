package user

import "time"

type User struct {
	ID        string
	TenantID  string
	Email     string
	Name      string
	CreatedAt time.Time
}

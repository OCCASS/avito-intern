package entity

type User struct {
	Id       string `json:"user_id" db:"id"`
	Name     string `json:"username" db:"name"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

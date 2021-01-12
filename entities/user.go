package entities

import "time"

type User struct {
	ID        int        `gorm:"primary_key" json:"id"`
	NickName  string     `db:"nick_name" json:"nick_name"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

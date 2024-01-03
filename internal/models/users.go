package models

import "time"

type User struct {
	ID        string     `db:"id" json:"id"`
	NickName  string     `db:"nick_name" json:"nick_name"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"password"`
	PhotoURL  string     `db:"photo_url" json:"photo_url"`
	Salt      string     `db:"salt" json:"salt"`
	Role      string     `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (user User) ToResponse() UserResponse {
	return UserResponse{
		ID:        user.ID,
		NickName:  user.NickName,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		PhotoURL:  user.PhotoURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type UserRequest struct {
	NickName string `json:"nick_name" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	NickName  string    `json:"nick_name"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	PhotoURL  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthToken struct {
	Token string `json:"token"`
}

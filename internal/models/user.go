package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Fullname  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}



// HashPassword hashes a plain-text password

func(u *User) HashPassword(password string) error{
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil{
		return err
	}
	u.Password = string(hashed)
	return  nil
}

// CheckPassword compares hash with plain passwprd

func (u *User) CheckPassword(password string) bool{
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	)
	return  err == nil
}

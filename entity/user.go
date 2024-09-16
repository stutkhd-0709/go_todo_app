package entity

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserID int64

type User struct {
	ID       UserID    `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Password string    `json:"password" db:"password"`
	Role     string    `json:"role" db:"role"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

func (u *User) ComparePassword(pw string) error {
	fmt.Println("aaaaaaaaaaaaaaaaaaaaa")
	fmt.Println(u.Password)
	fmt.Println(pw)
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
}

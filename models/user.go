package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"not null" json:"email"`
	FullName string `gorm:"size:255;not null" json:"fullName"`
	Password []byte `gorm:"not null" json:"password" validate:"required,min=8"`
}

func (user *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return nil
}
func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

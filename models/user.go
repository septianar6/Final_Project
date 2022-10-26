package models

import (
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username     string        `gorm:"uniqueIndex;not null" json:"username" form:"username" valid:"required~Your username is required"`
	Email        string        `gorm:"uniqueIndex;not null" json:"email" form:"email" valid:"required~Your email is required,email~Invalid email format"`
	Password     string        `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required,minstringlength(6)~passwrod hast to have a minimum length of 6 characters"`
	Age          int           `gorm:"not null" json:"age" form:"age"`
	Photos       []Photo       `json:"photos" gorm:"constraint:OnDelete:CASCADE;"`
	SocialMedias []SocialMedia `json:"social_medias" gorm:"constraint:OnDelete:CASCADE;"`
	Comments     []Comment     `json:"comments" gorm:"constraint:OnDelete:SET NULL;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	if !(u.Age > 8) {
		err = errors.New("Minimum Age is 9")
		return
	}

	err = nil
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(strings.TrimSpace(u.Username)) == 0 {
		err = errors.New("Username is required")
		return
	}
	if len(strings.TrimSpace(u.Email)) == 0 {
		err = errors.New("Email is required")
		return
	}

	if !(govalidator.IsEmail(u.Email)) {
		err = errors.New("Email format is invalid")
		return
	}

	err = nil
	return
}

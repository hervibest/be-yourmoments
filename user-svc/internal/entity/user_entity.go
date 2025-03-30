package entity

import (
	"database/sql"
	"time"
)

type User struct {
	Id                    string         `db:"id"`
	Username              string         `db:"username"`
	Email                 sql.NullString `db:"email"`
	EmailVerifiedAt       *time.Time     `db:"email_verified_at"`
	Password              sql.NullString `db:"password"`
	PhoneNumber           sql.NullString `db:"phone_number"`
	PhoneNumberVerifiedAt *time.Time     `db:"phone_number_verified_at"`
	GoogleId              sql.NullString `db:"google_id"`
	CreatedAt             *time.Time     `db:"created_at"`
	UpdatedAt             *time.Time     `db:"updated_at"`
}

func (u *User) HasEmail() bool {
	return u.Email.Valid
}

func (u *User) HasVerifiedEmail() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) HasPhoneNumber() bool {
	return u.PhoneNumber.Valid
}

func (u *User) HasVerifiedPhoneNumber() bool {
	return u.PhoneNumberVerifiedAt != nil
}

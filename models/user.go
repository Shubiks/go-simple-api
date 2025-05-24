package models

type User struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	// password is stored in DB but hidden in JSON
	Password string `db:"password" json:"-"`
}

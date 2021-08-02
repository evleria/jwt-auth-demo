// Package entities contains structs that reflect database entities
package entities

// User contains all data related to user and stored in database
type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	PassHash  string `db:"pass_hash"`
}

package auth

type User struct {
	Id        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	PassHash  string `db:"pass_hash"`
}

package auth

import "github.com/evleria/jwt-auth-demo/pkg/common/database"

type Repository interface {
	CreateNewUser(firstName, lastName, email, hash string) error
	GetUserByEmail(email string) (*User, error)
}

type repository struct {
	db database.Database
}

func NewRepository(database database.Database) Repository {
	return &repository{
		db: database,
	}
}

func (r *repository) CreateNewUser(firstName, lastName, email, hash string) error {
	_, err := r.db.Exec("INSERT INTO users (first_name, last_name, email, pass_hash) VALUES ($1, $2, $3, $4)",
		firstName, lastName, email, hash)
	return err
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	user := new(User)
	row := r.db.QueryRow("SELECT id, first_name, last_name, email, pass_hash FROM users WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)
	return user, err
}

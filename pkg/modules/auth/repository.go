package auth

import "github.com/evleria/jwt-auth-demo/pkg/common/database"

type Repository interface {
	CreateNewUser(firstName, lastName, email, hash string) error
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

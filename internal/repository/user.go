package repository

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/evleria/jwt-auth-demo/internal/repository/entities"
)

// User contains methods of create new user
type User interface {
	CreateNewUser(ctx context.Context, firstName, lastName, email, hash string) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, id int) (*entities.User, error)
}

type user struct {
	db *pgx.Conn
}

// NewUserRepository creates user repository
func NewUserRepository(db *pgx.Conn) User {
	return &user{
		db: db,
	}
}

func (r *user) CreateNewUser(ctx context.Context, firstName, lastName, email, hash string) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (first_name, last_name, email, pass_hash) VALUES ($1, $2, $3, $4)",
		firstName, lastName, email, hash)
	return err
}

func (r *user) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := new(entities.User)
	row := r.db.QueryRow(ctx, "SELECT id, first_name, last_name, email, pass_hash FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)
	return user, err
}

func (r *user) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	user := new(entities.User)
	row := r.db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)

	return user, err
}

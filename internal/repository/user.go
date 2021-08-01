package repository

import (
	"context"
	"github.com/evleria/jwt-auth-demo/internal/repository/entities"
	"github.com/jackc/pgx/v4"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, firstName, lastName, email, hash string) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserById(ctx context.Context, id int) (*entities.User, error)
}

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateNewUser(ctx context.Context, firstName, lastName, email, hash string) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (first_name, last_name, email, pass_hash) VALUES ($1, $2, $3, $4)",
		firstName, lastName, email, hash)
	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := new(entities.User)
	row := r.db.QueryRow(ctx, "SELECT id, first_name, last_name, email, pass_hash FROM users WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)
	return user, err
}

func (r *userRepository) GetUserById(ctx context.Context, id int) (*entities.User, error) {
	user := new(entities.User)
	row := r.db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)

	return user, err
}

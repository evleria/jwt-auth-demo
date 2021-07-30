package repositories

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type UserRepository interface {
	CreateNewUser(firstName, lastName, email, hash string) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
}

type User struct {
	Id        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	PassHash  string `db:"pass_hash"`
}

type userRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateNewUser(firstName, lastName, email, hash string) error {
	_, err := r.db.Exec(context.TODO(), "INSERT INTO users (first_name, last_name, email, pass_hash) VALUES ($1, $2, $3, $4)",
		firstName, lastName, email, hash)
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*User, error) {
	user := new(User)
	row := r.db.QueryRow(context.TODO(), "SELECT id, first_name, last_name, email, pass_hash FROM users WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)
	return user, err
}

func (r *userRepository) GetUserById(id int) (*User, error) {
	user := new(User)
	row := r.db.QueryRow(context.TODO(), "SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)

	return user, err
}

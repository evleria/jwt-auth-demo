package auth

import "github.com/evleria/jwt-auth-demo/internal/common/database"

type (
	UserRepository interface {
		CreateNewUser(firstName, lastName, email, hash string) error
		GetUserByEmail(email string) (*User, error)
		GetUserById(id int) (*User, error)
	}
)

type userRepository struct {
	db database.Database
}

func NewUserRepository(database database.Database) UserRepository {
	return &userRepository{
		db: database,
	}
}

func (r *userRepository) CreateNewUser(firstName, lastName, email, hash string) error {
	_, err := r.db.Exec("INSERT INTO users (first_name, last_name, email, pass_hash) VALUES ($1, $2, $3, $4)",
		firstName, lastName, email, hash)
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*User, error) {
	user := new(User)
	row := r.db.QueryRow("SELECT id, first_name, last_name, email, pass_hash FROM users WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)
	return user, err
}

func (r *userRepository) GetUserById(id int) (*User, error) {
	user := new(User)
	row := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.PassHash)

	return user, err
}

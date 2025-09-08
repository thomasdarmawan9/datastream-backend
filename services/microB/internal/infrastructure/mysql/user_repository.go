package mysql

import (
	"database/sql"
	"fmt"

	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/domain"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *domain.User) error {
	query := `INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, user.Username, user.PasswordHash, user.Role)
	return err
}

func (r *userRepo) FindByUsername(username string) (*domain.User, error) {
	query := `SELECT id, username, password_hash, role, created_at FROM users WHERE username = ?`
	row := r.db.QueryRow(query, username)

	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		fmt.Printf(">>> Query failed for username=%s, err=%v\n", username, err)
		return nil, err
	}

	fmt.Printf(">>> DB Found User: %+v\n", u)
	return &u, nil
}

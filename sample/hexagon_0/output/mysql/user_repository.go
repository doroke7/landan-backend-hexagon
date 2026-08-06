package mysql

import (
	"database/sql"
	"errors"

	"example/domain"
	"example/output/port"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) port.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AddOne(user *domain.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, name) VALUES (?, ?)
		 ON DUPLICATE KEY UPDATE name = VALUES(name)`,
		user.ID, user.Name,
	)
	return err
}

func (r *UserRepository) ShowOneById(id int) (*domain.User, error) {
	row := r.db.QueryRow(`SELECT id, name FROM users WHERE id = ?`, id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &user, nil
}

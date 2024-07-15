package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Id             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}
type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stamt := `INSERT INTO users(name,email,hashed_password,created) VALUES ($1,$2,$3,NOW())`
	_, err = m.DB.Exec(stamt, name, email, string(hashedPassword))
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" && strings.Contains(pqErr.Message, "users_uc_email") {

			return ErrDuplicateEmail
		}
		return err
	}
	return nil

}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id,hashed_password FROM users WHERE email=$1"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}

	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil

}
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

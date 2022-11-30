package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	Username   string    `json:"username"`
	Password   []byte    `json:"password"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Role       string    `json:"role"`
	Created_at time.Time `json:"createdAt"`
}

// var MapUsers = map[string]User{}

type UserModel struct {
	DB *sql.DB
}

// func (u *UserModel) CreateAdmin()

// Retrieve one user given the username
func (u *UserModel) Get(username string) (*User, error) {
	statement := `SELECT * FROM users WHERE username = ?`

	// row is a pointer to sql.Row which holds result from the DB
	row := u.DB.QueryRow(statement, username)

	user := &User{}

	err := row.Scan(&user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Role, &user.Created_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("the chosen user cannot be found")
		} else {
			return nil, err
		}
	}

	return user, nil
}

// Insert user struct & user data into record
func (u *UserModel) Insert(username string, password []byte, firstname string, lastname string, role string) (int, error) {
	statement := `INSERT INTO users (username, password, first_name, last_name, role, created_at) VALUES (?, ?, ?, ?, ?, UTC_TIMESTAMP())`

	// execute sql statement
	result, err := u.DB.Exec(statement, username, password, firstname, lastname, role)
	if err != nil {
		return 0, err
	}

	// get ID of newly inserted record
	num, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(num), nil
}

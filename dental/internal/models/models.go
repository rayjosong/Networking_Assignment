package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	Username   string
	Password   []byte
	FirstName  string
	LastName   string
	Role       string
	Created_at time.Time
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

// func (u *UserModel) ShowAll() (*User, error) {
// 	return nil, nil
// }

type Appointment struct {
	patient   User      `json:"patient"`
	startTime time.Time `json:"startTime"`
	endTime   time.Time `json:"endTime"`
	dentist   User      `json:"dentist"`
}

// return all appointments
// func (a *Appointment) GetAll() (*Appointment, error) {

// }

// // insert appointments
// func (a *Appointment) Insert(patient User, dentist User, startTime time.Time, endTime time.Time) {

// }

// edit appointments

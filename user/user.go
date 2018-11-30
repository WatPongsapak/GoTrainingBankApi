package user

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type UserApiServiceImp struct {
	DB *sql.DB
}

func (s *UserApiServiceImp) AllUsers() ([]User, error) {
	rows, err := s.DB.Query("SELECT id, first_name, last_name FROM users")
	if err != nil {
		return nil, err
	}
	users := []User{} // set empty slice without nil
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *UserApiServiceImp) CreateUser(user *User) error {
	row := s.DB.QueryRow("INSERT INTO users (first_name, last_name) values ($1, $2, $3) RETURNING id", user.FirstName, user.LastName)
	if err := row.Scan(&user.ID); err != nil {
		return err
	}
	return nil
}

func (s *UserApiServiceImp) GetUserByID(id int) (*User, error) {
	stmt := "SELECT id, first_name, last_name FROM users WHERE id = $1"
	row := s.DB.QueryRow(stmt, id)
	var user User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserApiServiceImp) DeleteUser(id int) error {
	stmt := "DELETE FROM users WHERE id = $1"
	_, err := s.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserApiServiceImp) UpdateUser(id int, user *User) (*User, error) {
	stmt := "UPDATE users SET first_name = $2,last_name = $2 WHERE id = $1"
	_, err := s.DB.Exec(stmt, id, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}
	return s.GetUserByID(id)
}

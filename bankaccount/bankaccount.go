package bankaccount

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Account struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"user_id"`
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	Balance       int64  `json:"balance"`
}

type AccountApiServiceImp struct {
	db *sql.DB
}

func (s *AccountApiServiceImp) AllAccounts() ([]Account, error) {
	rows, err := s.db.Query("SELECT id, first_name, last_name FROM accounts")
	if err != nil {
		return nil, err
	}
	accounts := []Account{} // set empty slice without nil
	for rows.Next() {
		var user Account
		err := rows.Scan(&user.ID, &user.UserID, &user.AccountNumber, &user.Name, &user.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, user)
	}
	return accounts, nil
}

func (s *AccountApiServiceImp) CreateAccount(user *Account) error {
	row := s.db.QueryRow("INSERT INTO accounts (first_name, last_name) values ($1, $2, $3) RETURNING id", user.FirstName, user.LastName)
	if err := row.Scan(&user.ID); err != nil {
		return err
	}
	return nil
}

func (s *AccountApiServiceImp) GetAccountByID(id int) (*Account, error) {
	stmt := "SELECT id, first_name, last_name FROM accounts WHERE id = $1"
	row := s.db.QueryRow(stmt, id)
	var user Account
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AccountApiServiceImp) DeleteAccount(id int) error {
	stmt := "DELETE FROM accounts WHERE id = $1"
	_, err := s.db.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountApiServiceImp) UpdateAccount(id int, user *Account) (*Account, error) {
	stmt := "UPDATE accounts SET first_name = $2,last_name = $2 WHERE id = $1"
	_, err := s.db.Exec(stmt, id, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}
	return s.GetAccountByID(id)
}

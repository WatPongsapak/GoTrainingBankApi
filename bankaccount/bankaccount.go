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
	DB *sql.DB
}

func (s *AccountApiServiceImp) AllAccounts() ([]Account, error) {
	rows, err := s.DB.Query("SELECT id, user_id, account_number, name, balance FROM bankaccounts")
	if err != nil {
		return nil, err
	}
	accounts := []Account{} // set empty slice without nil
	for rows.Next() {
		var acc Account
		err := rows.Scan(&acc.ID, &acc.UserID, &acc.AccountNumber, &acc.Name, &acc.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *AccountApiServiceImp) CreateAccount(acc *Account) error {
	row := s.DB.QueryRow("INSERT INTO bankaccounts (user_id, account_number, name) values ($1, $2, $3) RETURNING id", acc.UserID, acc.AccountNumber, acc.Name)
	if err := row.Scan(&acc.ID); err != nil {
		return err
	}
	return nil
}

func (s *AccountApiServiceImp) GetAccountByID(id int) (*Account, error) {
	stmt := "SELECT id, user_id, account_number, name, balance FROM bankaccounts WHERE id = $1"
	row := s.DB.QueryRow(stmt, id)
	var acc Account
	err := row.Scan(&acc.ID, &acc.UserID, &acc.AccountNumber, &acc.Name, &acc.Balance)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (s *AccountApiServiceImp) DeleteAccount(id int) error {
	stmt := "DELETE FROM accounts WHERE id = $1"
	_, err := s.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountApiServiceImp) UpdateAccount(id int, acc *Account) (*Account, error) {
	stmt := "UPDATE bankaccounts SET account_number = $2, name = $3, balance = $4  WHERE id = $1"
	_, err := s.DB.Exec(stmt, id, &acc.AccountNumber, &acc.Name, &acc.Balance)
	if err != nil {
		return nil, err
	}
	return s.GetAccountByID(id)
}

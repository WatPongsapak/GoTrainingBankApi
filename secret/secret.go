package secret

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Secret struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

type SecretServiceImp struct {
	db *sql.DB
}

func (s *SecretServiceImp) Insert(secret *Secret) error {
	row := s.db.QueryRow("INSERT INTO secrets (key) values ($1) RETURNING id", secret.Key)

	if err := row.Scan(&secret.ID); err != nil {
		return err
	}
	return nil
}

func (s *SecretServiceImp) FindSecretKey(secret *Secret) error {
	_, err := s.db.Query("SELECT key FROM secrets WHERE key=$1", secret.Key)
	if err != nil {
		return err
	}
	return nil
}

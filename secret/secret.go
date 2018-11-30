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
	DB *sql.DB
}

func (s *SecretServiceImp) Insert(secret *Secret) error {
	row := s.DB.QueryRow("INSERT INTO keys (key) values ($1) RETURNING id", secret.Key)

	if err := row.Scan(&secret.ID); err != nil {
		return err
	}
	return nil
}

func (s *SecretServiceImp) FindSecretKey(secret *Secret) error {
	_, err := s.DB.Query("SELECT key FROM keys WHERE key=$1", secret.Key)
	if err != nil {
		return err
	}
	return nil
}

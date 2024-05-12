package postresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4" //init postgresql drive
)

type Storage struct {
	db *pgx.Conn
}

func New(db_address, user, password, db_name string) (*Storage, error) {
	const fn = "storage.postgresql.New"

	config := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, db_address, db_name)
	//config := "postgres" + "://" + user + password + "@" + db_address + "/" + db_name

	connection, err := pgx.ParseConfig(config)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	db, err := pgx.ConnectConfig(context.Background(), connection)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = db.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS url(
		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const fn = "storage.postgresql.SaveURL"

	_, err := s.db.Exec(context.Background(), "INSERT INTO url(url, alias) VALUES ($1, $2)", urlToSave, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil

}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.postgresql.GetURL"

	var resUrl string
	err := s.db.QueryRow(context.Background(), "SELECT url FROM url WHERE alias = $1", alias).Scan(&resUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return " ", fmt.Errorf("%s: %w", fn, err)
	}
	if err != nil {
		return " ", fmt.Errorf("%s: %w", fn, err)
	}
	return resUrl, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const fn = "storage.postgresql.DeleteURL"

	_, err := s.db.Exec(context.Background(), "DELETE FROM url WHERE alias = $1", alias)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil

}

func (s *Storage) Close() error {
	err := s.db.Close(context.Background())
	return err
}

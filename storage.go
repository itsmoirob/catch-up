package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateMedia(*Media) error
	DeleteMedia(int) error
	UpdateMedia(*Media) error
	GetMedias() ([]*Media, error)
	GetMediaById(int) (*Media, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	//  docker run --name catch-up -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 -v postgres_volume:/var/lib/postgresql/data postgres
	connStr := "user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateMediaTable()
}

func (s *PostgresStore) CreateMediaTable() error {
	query := `CREATE TABLE IF NOT EXISTS media (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    form VARCHAR(255) NOT NULL,
    is_watched BOOLEAN NOT NULL,
    date_watched TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateMedia(media *Media) error {
	sqlStatement := `INSERT INTO media (title, form, is_watched, date_watched, created_at )
	VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(sqlStatement, media.Title, media.Form, media.IsWatched, media.DateWatched, media.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateMedia(*Media) error {
	return nil
}

func (s *PostgresStore) DeleteMedia(id int) error {
	_, err := s.db.Query("DELETE FROM media WHERE id = $1", id)
	return err
}

func (s *PostgresStore) GetMediaById(id int) (*Media, error) {
	rows, err := s.db.Query("SELECT * FROM media WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoMedia(rows)
	}
	return nil, fmt.Errorf("media %d not found", id)
}

func (s *PostgresStore) GetMedias() ([]*Media, error) {
	rows, err := s.db.Query("SELECT * FROM media")
	if err != nil {
		return nil, err
	}

	medias := []*Media{}

	for rows.Next() {
		media, err := scanIntoMedia(rows)
		if err != nil {
			return nil, err
		}

		medias = append(medias, media)
	}

	return medias, nil
}

func scanIntoMedia(rows *sql.Rows) (*Media, error) {
	media := new(Media)
	err := rows.Scan(
		&media.ID,
		&media.Title,
		&media.Form,
		&media.IsWatched,
		&media.DateWatched,
		&media.CreatedAt)

	return media, err
}

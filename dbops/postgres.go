package dbops

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresRepositoryImpl struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) PostgresRepository {
	return &PostgresRepositoryImpl{db: db}
}

func (p *PostgresRepositoryImpl) Insert(query string, args ...any) (*sql.Rows, error) {
	stmt, err := p.db.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(context.Background(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	return rows, nil
}

func (p *PostgresRepositoryImpl) Update(query string, args ...any) (*sql.Rows, error) {
	stmt, err := p.db.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(context.Background(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}

	return rows, nil
}

func (p *PostgresRepositoryImpl) Fetch(query string, args ...any) (*sql.Rows, error) {
	stmt, err := p.db.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare fetch statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(context.Background(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute fetch query: %w", err)
	}

	return rows, nil
}

func (p *PostgresRepositoryImpl) Delete(query string, args ...any) (*sql.Rows, error) {
	stmt, err := p.db.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(context.Background(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute delete query: %w", err)
	}

	return rows, nil
}

func (p *PostgresRepositoryImpl) GetStatus() error {
	if err := p.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping Postgres: %w", err)
	}
	return nil
}

func (p *PostgresRepositoryImpl) GetDB() *sql.DB {
	return p.db
}

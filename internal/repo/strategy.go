package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Strategy struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Params    string    `db:"params"`
	CreatedAt time.Time `db:"created_at"`
}

type StrategyRepo struct {
	db *sql.DB
}

func NewStrategyRepo(db *sql.DB) *StrategyRepo {
	return &StrategyRepo{db: db}
}

func (r *StrategyRepo) Create(ctx context.Context, s *Strategy) (int64, error) {
	query := `
		INSERT INTO strategies (name, type, params)
		VALUES (?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query, s.Name, s.Type, s.Params)
	if err != nil {
		return 0, fmt.Errorf("failed to insert strategy: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func (r *StrategyRepo) GetAll(ctx context.Context) ([]Strategy, error) {
	query := `
		SELECT id, name, type, params, created_at 
		FROM strategies 
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select strategies: %w", err)
	}
	defer rows.Close()

	var strategies []Strategy
	for rows.Next() {
		var s Strategy
		err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Params, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan strategy row: %w", err)
		}
		strategies = append(strategies, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return strategies, nil
}

func (r *StrategyRepo) GetByID(ctx context.Context, id int64) (*Strategy, error) {
	query := `SELECT id, name, type, params, created_at FROM strategies WHERE id = ?`

	var s Strategy
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.Type, &s.Params, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy by id: %w", err)
	}

	return &s, nil
}

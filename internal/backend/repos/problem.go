package repos

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/velosypedno/jobshop/pkg/dag/core"
)

type ProblemRecord struct {
	ID              uuid.UUID
	FactoryConfigID uuid.UUID
	CreatedAt       time.Time
	Order           core.Order
	Problem         core.Problem
}

type ProblemRepo struct {
	db *sql.DB
}

func NewProblemRepo(db *sql.DB) *ProblemRepo {
	return &ProblemRepo{
		db: db,
	}
}

func (r *ProblemRepo) Create(ctx context.Context, factoryConfigID uuid.UUID, order core.Order, problem core.Problem) (uuid.UUID, error) {
	id := uuid.New()
	createdAt := time.Now()

	orderJSON, err := order.ToJSON()
	if err != nil {
		return uuid.Nil, err
	}

	problemJSON, err := problem.ToJSON()
	if err != nil {
		return uuid.Nil, err
	}

	query := `
		INSERT INTO problems (id, factory_config_id, created_at, "order", problem)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.ExecContext(ctx, query, id, factoryConfigID, createdAt, orderJSON, problemJSON)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *ProblemRepo) GetByID(ctx context.Context, id uuid.UUID) (ProblemRecord, error) {
	query := `
		SELECT id, factory_config_id, created_at, "order", problem
		FROM problems
		WHERE id = $1
	`

	var rec ProblemRecord
	var orderJSON, problemJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rec.ID,
		&rec.FactoryConfigID,
		&rec.CreatedAt,
		&orderJSON,
		&problemJSON,
	)
	if err != nil {
		return ProblemRecord{}, err
	}

	order, err := core.OrderFromJSON(orderJSON)
	if err != nil {
		return ProblemRecord{}, err
	}
	rec.Order = order

	prob, err := core.ProblemFromJSON(problemJSON)
	if err != nil {
		return ProblemRecord{}, err
	}
	rec.Problem = prob

	return rec, nil
}

func (r *ProblemRepo) List(ctx context.Context) ([]ProblemRecord, error) {
	query := `
		SELECT id, factory_config_id, created_at, "order", problem
		FROM problems
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ProblemRecord
	for rows.Next() {
		var rec ProblemRecord
		var orderJSON, problemJSON []byte

		if err := rows.Scan(&rec.ID, &rec.FactoryConfigID, &rec.CreatedAt, &orderJSON, &problemJSON); err != nil {
			return nil, err
		}

		order, err := core.OrderFromJSON(orderJSON)
		if err != nil {
			return nil, err
		}
		rec.Order = order

		prob, err := core.ProblemFromJSON(problemJSON)
		if err != nil {
			return nil, err
		}
		rec.Problem = prob

		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *ProblemRepo) ListByFactoryConfigID(ctx context.Context, factoryConfigID uuid.UUID) ([]ProblemRecord, error) {
	query := `
		SELECT id, factory_config_id, created_at, "order", problem
		FROM problems
		WHERE factory_config_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, factoryConfigID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ProblemRecord
	for rows.Next() {
		var rec ProblemRecord
		var orderJSON, problemJSON []byte

		if err := rows.Scan(&rec.ID, &rec.FactoryConfigID, &rec.CreatedAt, &orderJSON, &problemJSON); err != nil {
			return nil, err
		}

		order, err := core.OrderFromJSON(orderJSON)
		if err != nil {
			return nil, err
		}
		rec.Order = order

		prob, err := core.ProblemFromJSON(problemJSON)
		if err != nil {
			return nil, err
		}
		rec.Problem = prob

		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *ProblemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM problems WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

package repos

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/velosypedno/jobshop/pkg/dag/core"
)

type FactoryConfigRecord struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Version   string
	Name      *string
	Config    core.FactoryConfig
}

type FactoryConfigRepo struct {
	db *sql.DB
}

func NewFactoryConfigRepo(db *sql.DB) *FactoryConfigRepo {
	return &FactoryConfigRepo{
		db: db,
	}
}

func (r *FactoryConfigRepo) Create(ctx context.Context, name *string, config core.FactoryConfig) (uuid.UUID, error) {
	id := uuid.New()
	createdAt := time.Now()
	const version string = "dag-v1"

	jsonData, err := config.ToJSON()
	if err != nil {
		return uuid.Nil, err
	}

	query := `
		INSERT INTO factory_configurations (id, created_at, version, name, config)
		VALUES ($1, $2, $3, $4, $5)
	`

	var nameVal any = nil
	if name != nil {
		nameVal = *name
	}

	_, err = r.db.ExecContext(ctx, query, id, createdAt, version, nameVal, jsonData)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *FactoryConfigRepo) GetByID(ctx context.Context, id uuid.UUID) (FactoryConfigRecord, error) {
	query := `
		SELECT id, created_at, version, name, config
		FROM factory_configurations
		WHERE id = $1
	`

	var rec FactoryConfigRecord
	var jsonData []byte
	var nameNull sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rec.ID,
		&rec.CreatedAt,
		&rec.Version,
		&nameNull,
		&jsonData,
	)
	if err != nil {
		return FactoryConfigRecord{}, err
	}

	if nameNull.Valid {
		rec.Name = &nameNull.String
	}

	cfg, err := core.FromJSON(jsonData)
	if err != nil {
		return FactoryConfigRecord{}, err
	}
	rec.Config = cfg

	return rec, nil
}

func (r *FactoryConfigRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM factory_configurations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *FactoryConfigRepo) List(ctx context.Context) ([]FactoryConfigRecord, error) {
	query := `
		SELECT id, created_at, version, name, config
		FROM factory_configurations
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []FactoryConfigRecord
	for rows.Next() {
		var rec FactoryConfigRecord
		var jsonData []byte
		var nameNull sql.NullString

		if err := rows.Scan(&rec.ID, &rec.CreatedAt, &rec.Version, &nameNull, &jsonData); err != nil {
			return nil, err
		}

		if nameNull.Valid {
			rec.Name = &nameNull.String
		}

		cfg, err := core.FromJSON(jsonData)
		if err != nil {
			return nil, err
		}
		rec.Config = cfg

		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

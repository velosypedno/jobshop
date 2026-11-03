CREATE TABLE factory_configurations (
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    version    VARCHAR(64) NOT NULL,
    name       VARCHAR(255),
    config     JSONB NOT NULL
);

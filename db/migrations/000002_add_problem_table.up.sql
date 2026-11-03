CREATE TABLE problems (
    id                UUID PRIMARY KEY,
    factory_config_id UUID NOT NULL REFERENCES factory_configurations(id),
    created_at        TIMESTAMP NOT NULL,
    "order"           JSONB NOT NULL,
    problem           JSONB NOT NULL
);

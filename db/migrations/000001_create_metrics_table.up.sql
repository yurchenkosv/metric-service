CREATE TABLE IF NOT EXISTS  metrics(
    id serial PRIMARY KEY,
    metric_id VARCHAR(256) NOT NULL UNIQUE,
    metric_type VARCHAR(50),
    metric_delta INTEGER,
    metric_value DOUBLE PRECISION,
    hash VARCHAR(300)
);
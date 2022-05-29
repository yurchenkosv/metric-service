CREATE TABLE IF NOT EXISTS  metrics(
    id serial PRIMARY KEY,
    metric_id varchar(256) NOT NULL,
    metric_delta INTEGER,
    metric_value DOUBLE PRECISION,
    hash varchar(300)
);
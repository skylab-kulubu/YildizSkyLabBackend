CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    data BYTEA NOT NULL,
    url TEXT NOT NULL,
    created_by SERIAL NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
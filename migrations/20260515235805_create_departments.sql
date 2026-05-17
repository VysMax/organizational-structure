-- +goose Up
CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE CASCADE
);

CREATE INDEX idx_departments_parent_id ON departments(parent_id);

-- +goose Down
DROP INDEX IF EXISTS idx_departments_parent_id;
DROP TABLE IF EXISTS departments;

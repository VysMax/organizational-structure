-- +goose Up
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    department_id INT NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);

CREATE INDEX idx_employees_department_id ON employees(department_id);

-- +goose Down
DROP INDEX IF EXISTS idx_employees_department_id;
DROP TABLE IF EXISTS employees;

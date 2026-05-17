-- +goose Up
CREATE UNIQUE INDEX idx_unique_child_name
ON departments (parent_id, name)
WHERE parent_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_unique_child_name;

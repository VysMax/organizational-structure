-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_circular_reference()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.parent_id IS NOT NULL THEN
        IF EXISTS (
            WITH RECURSIVE ancestors AS (
                SELECT id, parent_id FROM departments WHERE id = NEW.parent_id
                UNION ALL
                SELECT d.id, d.parent_id FROM departments d
                INNER JOIN ancestors a ON d.id = a.parent_id
            )
            SELECT 1 FROM ancestors WHERE id = NEW.id
        ) THEN
            RAISE EXCEPTION 'Circular reference detected';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_circular_reference
    BEFORE INSERT OR UPDATE ON departments
    FOR EACH ROW
    EXECUTE FUNCTION check_circular_reference();

-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS prevent_circular_reference ON departments;
DROP FUNCTION IF EXISTS check_circular_reference();

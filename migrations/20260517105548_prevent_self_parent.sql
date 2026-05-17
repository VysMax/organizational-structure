-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_self_parent()
RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.parent_id IS NOT NULL AND NEW.parent_id = NEW.id THEN
        RAISE EXCEPTION 'Department cannot be its own child';
    END IF;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER prevent_self_parent
    BEFORE INSERT OR UPDATE ON departments
    FOR EACH ROW
    EXECUTE FUNCTION check_self_parent();

-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS prevent_self_parent ON departments;
DROP FUNCTION IF EXISTS check_self_parent();

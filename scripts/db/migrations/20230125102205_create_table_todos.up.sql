
-- File created by: ./bin/db-tool new create_table_todos
BEGIN;
-- your migration here

-- See: https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-char-varchar-text/

-- NOTE https://samu.space/uuids-with-postgres-and-gorm/
--      thanks @anschnei
--      Consider to use UUID as the primary key
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    uuid UUID UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    due_date TIMESTAMP DEFAULT NULL
);

COMMIT;

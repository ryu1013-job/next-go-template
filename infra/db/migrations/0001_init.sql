-- +migrate Up
CREATE TABLE IF NOT EXISTS todos (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL CHECK (length(title) BETWEEN 1 AND 200),
  description TEXT,
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','done','archived')),
  due_date TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
-- +migrate Down
DROP TABLE IF EXISTS todos;

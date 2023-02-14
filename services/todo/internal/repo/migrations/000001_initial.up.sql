CREATE TABLE lists (
  id SERIAL PRIMARY KEY,
  uniq_id TEXT UNIQUE NOT NULL,

  title TEXT NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,

  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TYPE task_status AS ENUM ('completed', 'pending');

CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  uniq_id TEXT UNIQUE NOT NULL,
  task_uniq_id TEXT NOT NULL REFERENCES lists(uniq_id),

  title TEXT NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,
  status task_status NOT NULL,

  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
)

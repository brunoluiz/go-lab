CREATE TABLE radars (
  id SERIAL PRIMARY KEY,
  uniq_id TEXT UNIQUE NOT NULL,

  title TEXT NOT NULL,

  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE radar_quadrants (
  id SERIAL PRIMARY KEY,
  radar_id TEXT NOT NULL REFERENCES radars(id),

  name TEXT NOT NULL
);

CREATE TABLE radar_items (
  id SERIAL PRIMARY KEY,
  uniq_id TEXT UNIQUE NOT NULL,
  radar_id INTEGER NOT NULL REFERENCES radars(id),
  quadrant_id INTEGER NOT NULL REFERENCES radar_quadrants(id),

  title TEXT NOT NULL,

  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

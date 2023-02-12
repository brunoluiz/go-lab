create table lists (
  id SERIAL PRIMARY KEY,
  uid text not null,
  title text not null,
  created_at timestamp not null default now()
);


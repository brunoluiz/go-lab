create table lists (
  id text PRIMARY KEY,
  title text not null,
  created_at timestamp default now()
);


CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  phone_number TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  full_name TEXT,
  created_at TIMESTAMP NOT NULL
);

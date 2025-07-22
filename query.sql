CREATE TABLE users (
  usertag VARCHAR(9) PRIMARY KEY,
  firstname VARCHAR(255),
  lastname VARCHAR(255),
  email VARCHAR(255),
  telephone VARCHAR(12),
  gender VARCHAR(9),
  dob VARCHAR(255),
  password TEXT,
  photo_url TEXT,
  otp VARCHAR(6),
  otp_expiry TIMESTAMP
);

CREATE TYPE user_role AS ENUM ('admin', 'doctor', 'user');

ALTER TABLE users ADD COLUMN role user_role DEFAULT 'user';

CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,
    patient VARCHAR REFERENCES users(usertag),
    doctor VARCHAR REFERENCES users(usertag),
    scheduled_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    users VARCHAR REFERENCES users(usertag),
    item_name VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  users VARCHAR REFERENCES users(usertag),
  amount DECIMAL(20, 2) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'completed',
  payment_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

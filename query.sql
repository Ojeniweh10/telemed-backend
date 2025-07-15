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
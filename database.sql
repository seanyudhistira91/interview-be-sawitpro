/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */


CREATE TABLE users (
  id serial PRIMARY KEY,
  phone_number VARCHAR(13) UNIQUE NOT NULL,
  full_name VARCHAR(60) NOT NULL,
  hash_password VARCHAR UNIQUE NOT NULL,
  count_login BIGINT DEFAULT 0,
  created_at TIMESTAMP default CURRENT_TIMESTAMP,
  updated_at TIMESTAMP default CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE  FUNCTION update_updated_at_users()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON
        users
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_users();
CREATE TABLE IF NOT EXISTS users {
    id bigserial PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    password bytea NOT NULL
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
}
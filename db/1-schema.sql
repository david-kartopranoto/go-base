CREATE TABLE IF NOT EXISTS app_user (
    id BIGSERIAL,
    email varchar(500) NOT NULL,
    password varchar(1000) NOT NULL,
    username varchar(500) NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (id)
)
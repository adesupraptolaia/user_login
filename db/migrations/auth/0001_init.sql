-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    ksuid VARCHAR(255) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULl,
    PRIMARY KEY(ksuid)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

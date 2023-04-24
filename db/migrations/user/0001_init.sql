-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_profiles (
    user_ksuid VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    date_of_birth DATE,
    address VARCHAR(255),
    PRIMARY KEY(user_ksuid)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE user_profiles;

-- +goose StatementEnd
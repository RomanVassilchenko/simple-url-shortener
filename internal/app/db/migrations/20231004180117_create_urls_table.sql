-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
                                    id SERIAL PRIMARY KEY,
                                    alias VARCHAR(255) NOT NULL UNIQUE,
                                    url VARCHAR(255) NOT NULL
--                                     user_id INTEGER NOT NULL -- По рекомендации в лекции не делаю прямую зависимость
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd

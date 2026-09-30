-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Rates (
    rate_id SERIAL PRIMARY KEY,
    cur_from Text,
    cur_to Text,
    rate Float4
);

INSERT INTO Rates (cur_from, cur_to, rate) VALUES
    ('USD', 'RUB', 90.00),
    ('RUB', 'USD', 0.011111),
    ('EUR', 'RUB', 98.00),
    ('RUB', 'EUR', 0.010204),
    ('USD', 'EUR', 0.918367),
    ('EUR', 'USD', 1.088889);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS Rates;
-- +goose StatementEnd

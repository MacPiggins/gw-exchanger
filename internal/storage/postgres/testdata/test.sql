CREATE TABLE IF NOT EXISTS Rates (
    rate_id SERIAL PRIMARY KEY,
    cur_from Text,
    cur_to Text,
    rate Float4
);

INSERT INTO Rates (cur_from, cur_to, rate)
VALUES
	('USD', 'EUR', 0.92),
	('EUR', 'USD', 1.08);

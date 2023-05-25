CREATE TYPE currency AS ENUM ('UAH', 'USD', 'EUR');

CREATE TABLE IF NOT EXISTS currency_rate (
    rate_date date PRIMARY KEY,
    base_currency currency,
    target_currency currency,
    sale numeric(9, 6), -- 100.000_001
    purchase numeric(9, 6), -- 100.000_001
    source VARCHAR(7)
);

CREATE INDEX idx_currency_rate_base_currency ON currency_rate (base_currency);
CREATE INDEX idx_currency_rate_target_currency ON currency_rate (target_currency);


---- create above / drop below ----

DROP TYPE currency;

DROP TABLE currency_rate;

CREATE TYPE currency AS ENUM ('UAH', 'USD', 'EUR');
CREATE TYPE exchange_source AS ENUM ('NBU', 'Privat');

CREATE TABLE IF NOT EXISTS currency_rate (
    id Bigserial PRIMARY KEY NOT NULL,
    rate_date date,
    base_currency currency,
    target_currency currency,
    sale numeric(9, 6), -- 100.000_001
    purchase numeric(9, 6), -- 100.000_001
    source exchange_source
);

CREATE INDEX idx_currency_rate_base_currency ON currency_rate (base_currency);
CREATE INDEX idx_currency_rate_target_currency ON currency_rate (target_currency);


---- create above / drop below ----

DROP TYPE currency;

DROP TABLE currency_rate;

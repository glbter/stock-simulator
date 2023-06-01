create extension if not exists "uuid-ossp";

CREATE TABLE IF NOT EXISTS ticker (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(7) NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS split (
    date date NOT NULL DEFAULT date(now()),
    ticker_id UUID REFERENCES ticker(id) ON DELETE CASCADE NOT NULL,
    before NUMERIC(6,4) NOT NULL, -- 99.9999
    after NUMERIC(6,4) NOT NULL, -- 99.9999

    PRIMARY KEY (ticker_id, date)
);

CREATE TABLE IF NOT EXISTS stock_daily (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticker_id UUID NOT NULL REFERENCES ticker(id) ON DELETE CASCADE,
    date date NOT NULL DEFAULT date(now()),
    high NUMERIC(8,3),
    low NUMERIC(8,4),
    open NUMERIC(8,3),
    close NUMERIC(8,3),
    volume NUMERIC(14,3)
);

CREATE TYPE portfolio_action AS ENUM ('BUY', 'SELL');


CREATE TABLE IF NOT EXISTS portfolio_record(
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    investor_id VARCHAR(36) NOT NULL,
    ticker_id UUID NOT NULL,
    date date NOT NULL DEFAULT date(now()),
    amount NUMERIC(7,4) NOT NULL, -- 999.9999
    price_usd NUMERIC(9,4) NOT NULL, -- 99999.9999
    action portfolio_action NOT NULL
);

CREATE AGGREGATE mul(numeric) (INITCOND = 1, STYPE = numeric, SFUNC = numeric_mul);


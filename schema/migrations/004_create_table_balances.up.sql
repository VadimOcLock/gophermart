CREATE TABLE balances
(
    id                SERIAL PRIMARY KEY,
    user_id           INT            NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    current_balance   NUMERIC(10, 2) NOT NULL DEFAULT 0,
    withdrawn_balance NUMERIC(10, 2) NOT NULL DEFAULT 0
);

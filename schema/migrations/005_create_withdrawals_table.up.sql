CREATE TABLE withdrawals
(
    id           SERIAL PRIMARY KEY,
    user_id      INT            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number TEXT           NOT NULL,
    sum          NUMERIC(10, 2) NOT NULL CHECK ( sum > 0 ),
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_withdrawals_user_id ON withdrawals (user_id);
CREATE INDEX idx_withdrawals_processed_at ON withdrawals (processed_at DESC);
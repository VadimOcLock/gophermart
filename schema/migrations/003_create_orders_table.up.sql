CREATE TABLE orders
(
    id           SERIAL PRIMARY KEY,
    user_id      INT  NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    order_number TEXT NOT NULL UNIQUE,
    status       TEXT NOT NULL CHECK ( status in ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED') ),
    accrual      NUMERIC(10, 2) DEFAULT 0,
    uploaded_at  TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, order_number)
);

CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_orders_uploaded_at ON orders (uploaded_at DESC);
CREATE TABLE IF NOT EXISTS user_stock_holdings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE RESTRICT,
    total_quantity NUMERIC(18, 6) NOT NULL DEFAULT 0,
    average_price NUMERIC(18, 4) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, stock_id)
);

CREATE INDEX IF NOT EXISTS idx_user_stock_holdings_user_id ON user_stock_holdings(user_id);
CREATE INDEX IF NOT EXISTS idx_user_stock_holdings_stock_id ON user_stock_holdings(stock_id);

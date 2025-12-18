CREATE TABLE IF NOT EXISTS reward_events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(id) ON DELETE RESTRICT,
    quantity NUMERIC(18, 6) NOT NULL CHECK (quantity > 0),
    stock_price NUMERIC(18, 4) NOT NULL,
    total_value NUMERIC(18, 4) NOT NULL,
    event_type VARCHAR(50) NOT NULL DEFAULT 'STOCK_REWARD',
    status VARCHAR(50) NOT NULL DEFAULT 'COMPLETED',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reward_events_user_id ON reward_events(user_id);
CREATE INDEX IF NOT EXISTS idx_reward_events_stock_id ON reward_events(stock_id);
CREATE INDEX IF NOT EXISTS idx_reward_events_created_at ON reward_events(created_at);
CREATE INDEX IF NOT EXISTS idx_reward_events_status ON reward_events(status);

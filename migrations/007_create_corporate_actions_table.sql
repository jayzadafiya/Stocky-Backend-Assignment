-- Create corporate_actions table
CREATE TABLE IF NOT EXISTS corporate_actions (
    id SERIAL PRIMARY KEY,
    stock_id INT NOT NULL REFERENCES stocks(id),
    action_type VARCHAR(20) NOT NULL CHECK (action_type IN ('STOCK_SPLIT', 'MERGER', 'DELISTING')),
    split_ratio NUMERIC(10, 4),
    merger_to_stock_id INT REFERENCES stocks(id),
    merger_ratio NUMERIC(10, 4),
    effective_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED')),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_corporate_actions_stock_id ON corporate_actions(stock_id);
CREATE INDEX IF NOT EXISTS idx_corporate_actions_status ON corporate_actions(status);
CREATE INDEX IF NOT EXISTS idx_corporate_actions_effective_date ON corporate_actions(effective_date);

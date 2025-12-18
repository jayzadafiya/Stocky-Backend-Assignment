CREATE TABLE IF NOT EXISTS ledger_entries (
    id SERIAL PRIMARY KEY,
    reward_event_id INTEGER NOT NULL REFERENCES reward_events(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entry_type VARCHAR(50) NOT NULL, 
    account_type VARCHAR(50) NOT NULL, 
    stock_id INTEGER REFERENCES stocks(id) ON DELETE RESTRICT,
    quantity NUMERIC(18, 6),
    amount NUMERIC(18, 4), 
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ledger_entries_reward_event_id ON ledger_entries(reward_event_id);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_user_id ON ledger_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_entry_type ON ledger_entries(entry_type);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_account_type ON ledger_entries(account_type);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_stock_id ON ledger_entries(stock_id);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_created_at ON ledger_entries(created_at);

ALTER TABLE ledger_entries ADD CONSTRAINT check_quantity_or_amount 
    CHECK ((quantity IS NOT NULL AND amount IS NULL) OR (quantity IS NULL AND amount IS NOT NULL));

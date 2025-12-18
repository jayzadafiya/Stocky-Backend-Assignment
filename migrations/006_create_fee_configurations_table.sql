CREATE TABLE IF NOT EXISTS fee_configurations (
    id SERIAL PRIMARY KEY,
    fee_type VARCHAR(50) UNIQUE NOT NULL,
    percentage NUMERIC(5, 4) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO fee_configurations (fee_type, percentage, description) VALUES
('BROKERAGE', 0.0010, 'Brokerage fee - 0.10% of transaction value'),
('STT', 0.0010, 'Securities Transaction Tax - 0.10% of transaction value'),
('GST', 0.1800, 'GST on brokerage - 18% of brokerage amount')
ON CONFLICT (fee_type) DO NOTHING;

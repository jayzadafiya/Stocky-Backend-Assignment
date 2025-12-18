# Database Schema Documentation

## Entity Relationship Diagram

```
┌─────────────────┐
│     USERS       │
├─────────────────┤
│ PK id           │───┐
│    email        │   │
│    name         │   │
│    phone        │   │
│    is_active    │   │
│    created_at   │   │
│    updated_at   │   │
└─────────────────┘   │
                      │
                      │
┌─────────────────┐   │         ┌──────────────────────┐
│     STOCKS      │   │         │   REWARD_EVENTS      │
├─────────────────┤   │         ├──────────────────────┤
│ PK id           │───┼────────→│ FK user_id           │
│    symbol       │   │         │ FK stock_id          │
│    name         │   │         │ PK id                │
│    exchange     │   │         │    quantity          │
│ current_price   │   │         │    stock_price       │
│    is_active    │   │         │    total_value       │
│    created_at   │   │         │    event_type        │
│    updated_at   │   │         │    status            │
└─────────────────┘   │         │    description       │
         │            │         │    created_at        │
         │            │         │    updated_at        │
         │            │         └──────────────────────┘
         │            │                     │
         │            │                     │
         │            │         ┌──────────────────────┐
         │            │         │  LEDGER_ENTRIES      │
         │            │         ├──────────────────────┤
         │            └────────→│ FK user_id           │
         │                      │ FK reward_event_id   │
         └─────────────────────→│ FK stock_id (null)   │
                                │ PK id                │
                                │    entry_type        │
                                │    account_type      │
                                │    quantity (null)   │
                                │    amount (null)     │
                                │    description       │
                                │    created_at        │
                                └──────────────────────┘
         │            │
         │            │         ┌──────────────────────────┐
         │            │         │ USER_STOCK_HOLDINGS      │
         │            │         ├──────────────────────────┤
         │            └────────→│ FK user_id               │
         └─────────────────────→│ FK stock_id              │
                                │ PK id                    │
                                │    total_quantity        │
                                │    average_price         │
                                │    created_at            │
                                │    updated_at            │
                                │ UQ (user_id, stock_id)   │
                                └──────────────────────────┘

         │
         │                      ┌──────────────────────────┐
         │                      │  CORPORATE_ACTIONS       │
         │                      ├──────────────────────────┤
         └─────────────────────→│ FK stock_id              │
                                │ FK merger_to_stock_id    │
                                │ PK id                    │
                                │    action_type           │
                                │    split_ratio           │
                                │    merger_ratio          │
                                │    effective_date        │
                                │    status                │
                                │    description           │
                                │    created_at            │
                                │    processed_at          │
                                └──────────────────────────┘

┌──────────────────────────┐
│  FEE_CONFIGURATIONS      │
├──────────────────────────┤
│ PK id                    │
│    fee_type              │
│    percentage            │
│    description           │
│    is_active             │
│    created_at            │
│    updated_at            │
└──────────────────────────┘
```

---

## Table Definitions

### 1. USERS

Stores user information.

| Column     | Type         | Constraints          | Description               |
| ---------- | ------------ | -------------------- | ------------------------- |
| id         | SERIAL       | PRIMARY KEY          | Auto-incrementing user ID |
| email      | VARCHAR(255) | UNIQUE, NOT NULL     | User email (unique)       |
| name       | VARCHAR(255) | NOT NULL             | User full name            |
| phone      | VARCHAR(20)  | NOT NULL             | User phone number         |
| is_active  | BOOLEAN      | DEFAULT true         | Active status             |
| created_at | TIMESTAMP    | DEFAULT CURRENT_TIME | Record creation time      |
| updated_at | TIMESTAMP    | DEFAULT CURRENT_TIME | Last update time          |

**Indexes:**

- Primary Key: `id`
- Unique: `email`

---

### 2. STOCKS

Stores stock/security information.

| Column        | Type          | Constraints          | Description                |
| ------------- | ------------- | -------------------- | -------------------------- |
| id            | SERIAL        | PRIMARY KEY          | Auto-incrementing stock ID |
| symbol        | VARCHAR(50)   | UNIQUE, NOT NULL     | Stock ticker symbol        |
| name          | VARCHAR(255)  | NOT NULL             | Company name               |
| exchange      | VARCHAR(50)   | NOT NULL             | Stock exchange (NSE, BSE)  |
| current_price | NUMERIC(18,4) | NOT NULL             | Current market price (₹)   |
| is_active     | BOOLEAN       | DEFAULT true         | Active/Delisted status     |
| created_at    | TIMESTAMP     | DEFAULT CURRENT_TIME | Record creation time       |
| updated_at    | TIMESTAMP     | DEFAULT CURRENT_TIME | Last update time           |

**Indexes:**

- Primary Key: `id`
- Unique: `symbol`

**Precision:** Prices stored with 4 decimal places for accuracy.

---

### 3. REWARD_EVENTS

Records all stock reward issuances and adjustments.

| Column      | Type          | Constraints          | Description                |
| ----------- | ------------- | -------------------- | -------------------------- |
| id          | SERIAL        | PRIMARY KEY          | Auto-incrementing event ID |
| user_id     | INTEGER       | FK → users(id)       | Recipient user             |
| stock_id    | INTEGER       | FK → stocks(id)      | Stock being rewarded       |
| quantity    | NUMERIC(18,6) | NOT NULL, > 0        | Number of shares           |
| stock_price | NUMERIC(18,4) | NOT NULL             | Price at time of reward    |
| total_value | NUMERIC(18,4) | NOT NULL             | quantity × stock_price     |
| event_type  | VARCHAR(50)   | DEFAULT 'REWARD'     | REWARD or ADJUSTMENT       |
| status      | VARCHAR(50)   | DEFAULT 'COMPLETED'  | Event status               |
| description | TEXT          |                      | Event description          |
| created_at  | TIMESTAMP     | DEFAULT CURRENT_TIME | Event timestamp            |
| updated_at  | TIMESTAMP     | DEFAULT CURRENT_TIME | Last update time           |

**Indexes:**

- Primary Key: `id`
- Foreign Keys: `user_id`, `stock_id`
- Index on: `user_id`, `created_at`, `event_type`

**Check Constraints:**

- `quantity > 0`

**Event Types:**

- `REWARD` - Normal stock reward issuance
- `ADJUSTMENT` - Refund/adjustment (reversal)

**Precision:** Quantity with 6 decimal places for fractional shares.

---

### 4. LEDGER_ENTRIES

Double-entry accounting ledger for all transactions.

| Column          | Type          | Constraints          | Description                |
| --------------- | ------------- | -------------------- | -------------------------- |
| id              | SERIAL        | PRIMARY KEY          | Auto-incrementing entry ID |
| reward_event_id | INTEGER       | FK → reward_events   | Related reward event       |
| user_id         | INTEGER       | FK → users(id)       | User account               |
| stock_id        | INTEGER       | FK → stocks(id)      | Stock (nullable)           |
| entry_type      | VARCHAR(50)   | NOT NULL             | DEBIT or CREDIT            |
| account_type    | VARCHAR(100)  | NOT NULL             | Account category           |
| quantity        | NUMERIC(18,6) |                      | Stock quantity (nullable)  |
| amount          | NUMERIC(18,4) |                      | INR amount (nullable)      |
| description     | TEXT          |                      | Entry description          |
| created_at      | TIMESTAMP     | DEFAULT CURRENT_TIME | Entry timestamp            |

**Indexes:**

- Primary Key: `id`
- Foreign Keys: `reward_event_id`, `user_id`, `stock_id`
- Index on: `user_id`, `account_type`, `created_at`

**Check Constraints:**

- `(quantity IS NOT NULL AND amount IS NULL) OR (quantity IS NULL AND amount IS NOT NULL)`
  - Either quantity OR amount must be set, never both

**Entry Types:**

- `DEBIT` - Increases asset accounts, decreases liability/equity
- `CREDIT` - Increases liability/equity, decreases assets

**Account Types:**

- `STOCK_UNITS` - Stock holdings (uses quantity)
- `INR_CASH` - Cash flow (uses amount)
- `BROKERAGE_FEE` - Brokerage charges (uses amount)
- `STT_FEE` - Securities Transaction Tax (uses amount)
- `GST_FEE` - GST on fees (uses amount)

**Double-Entry Example:**

```
Reward: 10 shares of RELIANCE @ ₹2450.75 = ₹24,507.50

DEBIT   STOCK_UNITS    10.000000 shares   (Asset increase)
CREDIT  INR_CASH       ₹24,507.50         (Cash outflow)
CREDIT  BROKERAGE_FEE  ₹12.25             (Fee expense)
CREDIT  STT_FEE        ₹24.51             (Tax expense)
CREDIT  GST_FEE        ₹2.21              (Tax expense)
```

---

### 5. USER_STOCK_HOLDINGS

Aggregated current stock holdings for each user.

| Column         | Type          | Constraints          | Description                  |
| -------------- | ------------- | -------------------- | ---------------------------- |
| id             | SERIAL        | PRIMARY KEY          | Auto-incrementing holding ID |
| user_id        | INTEGER       | FK → users(id)       | User account                 |
| stock_id       | INTEGER       | FK → stocks(id)      | Stock held                   |
| total_quantity | NUMERIC(18,6) | DEFAULT 0            | Total shares owned           |
| average_price  | NUMERIC(18,4) | DEFAULT 0            | Weighted average buy price   |
| created_at     | TIMESTAMP     | DEFAULT CURRENT_TIME | First holding timestamp      |
| updated_at     | TIMESTAMP     | DEFAULT CURRENT_TIME | Last update time             |

**Indexes:**

- Primary Key: `id`
- Unique: `(user_id, stock_id)`
- Foreign Keys: `user_id`, `stock_id`
- Index on: `user_id`, `total_quantity`

**Unique Constraint:** One row per user-stock combination.

**Average Price Calculation:**

```sql
new_avg = ((old_qty × old_avg) + (new_qty × new_price)) / (old_qty + new_qty)
```

**Note:** Holdings are updated via UPSERT (ON CONFLICT) operations.

---

### 6. CORPORATE_ACTIONS

Tracks stock splits, mergers, and delistings.

| Column             | Type          | Constraints          | Description                         |
| ------------------ | ------------- | -------------------- | ----------------------------------- |
| id                 | SERIAL        | PRIMARY KEY          | Auto-incrementing action ID         |
| stock_id           | INTEGER       | FK → stocks(id)      | Stock being affected                |
| action_type        | VARCHAR(50)   | NOT NULL             | Type of corporate action            |
| split_ratio        | NUMERIC(10,4) |                      | Split ratio (nullable)              |
| merger_to_stock_id | INTEGER       | FK → stocks(id)      | Target stock for merger (nullable)  |
| merger_ratio       | NUMERIC(10,4) |                      | Merger conversion ratio (nullable)  |
| effective_date     | DATE          | NOT NULL             | When action takes effect            |
| status             | VARCHAR(50)   | DEFAULT 'PENDING'    | PENDING or COMPLETED                |
| description        | TEXT          |                      | Action description                  |
| created_at         | TIMESTAMP     | DEFAULT CURRENT_TIME | Action creation time                |
| processed_at       | TIMESTAMP     |                      | When action was executed (nullable) |

**Indexes:**

- Primary Key: `id`
- Foreign Keys: `stock_id`, `merger_to_stock_id`
- Index on: `stock_id`, `status`, `effective_date`

**Check Constraints:**

- `action_type IN ('STOCK_SPLIT', 'MERGER', 'DELISTING')`
- `status IN ('PENDING', 'COMPLETED')`

**Action Types:**

1. **STOCK_SPLIT** (split_ratio required)

   - Multiplies all holdings by ratio
   - Divides stock price by ratio
   - Example: 1:2 split → ratio = 2.0

2. **MERGER** (merger_to_stock_id, merger_ratio required)

   - Converts holdings to target stock
   - Uses conversion ratio
   - Deactivates source stock
   - Example: A merges into B at 1:0.5 → ratio = 0.5

3. **DELISTING** (no additional fields)
   - Zeros all holdings
   - Deactivates stock
   - Blocks new rewards

**Processing Logic:**

- Actions are created as PENDING
- Must be explicitly processed via API
- Once processed, status changes to COMPLETED
- Cannot be processed twice

---

### 7. FEE_CONFIGURATIONS

Stores fee percentages for transaction costs.

| Column      | Type         | Constraints          | Description                   |
| ----------- | ------------ | -------------------- | ----------------------------- |
| id          | SERIAL       | PRIMARY KEY          | Auto-incrementing config ID   |
| fee_type    | VARCHAR(50)  | UNIQUE, NOT NULL     | Type of fee                   |
| percentage  | NUMERIC(5,4) | NOT NULL             | Fee percentage (e.g., 0.0005) |
| description | TEXT         |                      | Fee description               |
| is_active   | BOOLEAN      | DEFAULT true         | Active status                 |
| created_at  | TIMESTAMP    | DEFAULT CURRENT_TIME | Record creation time          |
| updated_at  | TIMESTAMP    | DEFAULT CURRENT_TIME | Last update time              |

**Indexes:**

- Primary Key: `id`
- Unique: `fee_type`

**Default Fee Types:**

- `BROKERAGE` - 0.0005 (0.05%)
- `STT` - 0.001 (0.1%)
- `GST` - 0.18 (18% on brokerage)

**Precision:** Up to 4 decimal places (0.0001 = 0.01%)

---

## Relationships

### One-to-Many

1. **users → reward_events**

   - One user can have many reward events
   - `users.id → reward_events.user_id`

2. **stocks → reward_events**

   - One stock can be in many reward events
   - `stocks.id → reward_events.stock_id`

3. **reward_events → ledger_entries**

   - One reward creates multiple ledger entries
   - `reward_events.id → ledger_entries.reward_event_id`

4. **users → ledger_entries**

   - One user has many ledger entries
   - `users.id → ledger_entries.user_id`

5. **stocks → ledger_entries**

   - One stock can appear in many ledger entries
   - `stocks.id → ledger_entries.stock_id` (nullable)

6. **users → user_stock_holdings**

   - One user can hold many stocks
   - `users.id → user_stock_holdings.user_id`

7. **stocks → user_stock_holdings**

   - One stock can be held by many users
   - `stocks.id → user_stock_holdings.stock_id`

8. **stocks → corporate_actions**
   - One stock can have many corporate actions
   - `stocks.id → corporate_actions.stock_id`

### Many-to-Many

1. **users ↔ stocks** (via user_stock_holdings)
   - Users can hold multiple stocks
   - Stocks can be held by multiple users
   - Junction table: `user_stock_holdings`
   - Unique constraint: `(user_id, stock_id)`

---

## Data Integrity Rules

### Cascading Actions

- **ON DELETE RESTRICT** for most foreign keys
  - Cannot delete users/stocks with existing records
  - Maintains data integrity and audit trail

### Check Constraints

1. **reward_events.quantity** > 0
2. **ledger_entries** - Either quantity OR amount (not both)
3. **corporate_actions.action_type** - Must be valid enum
4. **corporate_actions.status** - PENDING or COMPLETED

### Unique Constraints

1. **users.email** - One email per user
2. **stocks.symbol** - One symbol per stock
3. **user_stock_holdings(user_id, stock_id)** - One holding per user-stock pair
4. **fee_configurations.fee_type** - One config per fee type

---

## Indexing Strategy

### Primary Indexes

- All tables have `id` as PRIMARY KEY (B-tree)
- Sequential IDs for optimal INSERT performance

### Foreign Key Indexes

- All foreign keys automatically indexed
- Optimizes JOIN operations

### Query Optimization Indexes

```sql
-- Fast user reward lookups
CREATE INDEX idx_reward_events_user_created
ON reward_events(user_id, created_at DESC);

-- Fast stock holdings queries
CREATE INDEX idx_holdings_user_quantity
ON user_stock_holdings(user_id)
WHERE total_quantity > 0;

-- Fast corporate action status checks
CREATE INDEX idx_corporate_actions_stock_status
ON corporate_actions(stock_id, status);

-- Fast ledger entry lookups by account
CREATE INDEX idx_ledger_user_account
ON ledger_entries(user_id, account_type);
```

---

## Precision and Rounding

### Storage Precision

- **Stock Quantities:** NUMERIC(18, 6) - 6 decimal places
- **INR Amounts:** NUMERIC(18, 4) - 4 decimal places (stored)
- **Stock Prices:** NUMERIC(18, 4) - 4 decimal places
- **Fee Percentages:** NUMERIC(5, 4) - 4 decimal places

### Display Precision (API)

- **Stock Quantities:** 6 decimal places (as stored)
- **INR Amounts:** 2 decimal places (rounded via ROUND())
- **Stock Prices:** 2 decimal places (rounded via ROUND())

All INR calculations in queries use PostgreSQL's `ROUND(value, 2)` function to prevent floating-point rounding errors in the application layer.

---

## Sample Data Flow

### Creating a Reward

```
1. INSERT INTO reward_events
   - user_id: 1
   - stock_id: 1 (RELIANCE)
   - quantity: 10.5
   - stock_price: 2450.75
   - total_value: 25732.88

2. INSERT INTO ledger_entries (Stock Units)
   - reward_event_id: 1
   - account_type: STOCK_UNITS
   - quantity: 10.5
   - entry_type: DEBIT

3. INSERT INTO ledger_entries (Cash)
   - reward_event_id: 1
   - account_type: INR_CASH
   - amount: 25732.88
   - entry_type: CREDIT

4. INSERT INTO ledger_entries (Fees) × 3
   - BROKERAGE_FEE, STT_FEE, GST_FEE

5. UPSERT user_stock_holdings
   - Update if exists
   - Insert if new
   - Recalculate average_price
```

### Processing Stock Split

```
1. SELECT * FROM corporate_actions WHERE id = ?
   - Validate status = PENDING

2. UPDATE user_stock_holdings
   - total_quantity = total_quantity × split_ratio
   - average_price = average_price / split_ratio

3. UPDATE stocks
   - current_price = current_price / split_ratio

4. UPDATE corporate_actions
   - status = COMPLETED
   - processed_at = NOW()
```

---

## Database Maintenance

### Recommended Maintenance Tasks

1. **VACUUM ANALYZE** - Weekly on large tables
2. **REINDEX** - Monthly on high-write tables
3. **Partition** reward_events and ledger_entries by month/year for large datasets
4. **Archive** old COMPLETED corporate actions

### Backup Strategy

- Daily full backups
- Continuous WAL archiving
- Point-in-time recovery capability
- Test restore procedures monthly

# Stocky Backend Assignment - API Documentation

## Base URL

```
http://localhost:8080/api
```

## Table of Contents

- [Health Check](#health-check)
- [Reward Endpoints](#reward-endpoints)
- [User Endpoints](#user-endpoints)
- [Stock Endpoints](#stock-endpoints)
- [Corporate Action Endpoints](#corporate-action-endpoints)
- [Ledger Endpoints](#ledger-endpoints)

---

## Health Check

### GET /health

Check if the API is running.

**Response:**

```json
{
  "status": "ok",
  "message": "Stocky Backend Assignment API is running"
}
```

---

## Reward Endpoints

### 1. Create Reward

**POST** `/api/reward`

Issue stock rewards to users with double-entry ledger accounting.

**Request Body:**

```json
{
  "user_id": 1,
  "stock_symbol": "RELIANCE",
  "quantity": 10.5,
  "description": "Performance bonus Q4",
  "idempotency_key": "unique-key-123" // Optional
}
```

**Response:** `201 Created`

```json
{
  "message": "Reward created successfully with ledger entries",
  "data": {
    "id": 1,
    "user_id": 1,
    "stock_id": 1,
    "quantity": 10.5,
    "stock_price": 2450.75,
    "total_value": 25732.88,
    "event_type": "REWARD",
    "status": "COMPLETED",
    "description": "Performance bonus Q4",
    "created_at": "2025-12-19T10:30:00Z",
    "updated_at": "2025-12-19T10:30:00Z"
  }
}
```

**Validations:**

- User must exist and be active
- Stock must exist and be active (not delisted)
- No pending corporate actions on the stock
- Duplicate detection (5-minute window)
- Idempotency key validation (1-hour window)

**Error Responses:**

```json
{
  "error": "stock 'XYZ' is delisted and cannot receive new rewards"
}
```

```json
{
  "error": "stock 'RELIANCE' has a pending STOCK_SPLIT corporate action. Please process it before issuing new rewards"
}
```

```json
{
  "error": "duplicate reward detected: similar reward was created within the last 5 minutes"
}
```

---

### 2. Adjust/Refund Reward

**POST** `/api/reward/adjust`

Refund or partially refund previously issued rewards.

**Request Body:**

```json
{
  "reward_event_id": 123,
  "adjustment_type": "REFUND", // or "PARTIAL_REFUND"
  "quantity": 10.5,
  "reason": "Reward issued in error"
}
```

**Response:** `201 Created`

```json
{
  "message": "Reward adjusted successfully",
  "data": {
    "id": 124,
    "user_id": 1,
    "stock_id": 1,
    "quantity": 10.5,
    "stock_price": 2450.75,
    "total_value": 25732.88,
    "event_type": "ADJUSTMENT",
    "status": "COMPLETED",
    "description": "REFUND for reward #123: Reward issued in error",
    "created_at": "2025-12-19T11:00:00Z",
    "updated_at": "2025-12-19T11:00:00Z"
  }
}
```

**Validations:**

- Original reward event must exist
- Cannot adjust an adjustment
- REFUND must match original quantity
- PARTIAL_REFUND must be less than original
- User must have sufficient holdings

**Error Responses:**

```json
{
  "error": "reward event not found"
}
```

```json
{
  "error": "insufficient holdings: user has 5.000000, adjustment requires 10.500000"
}
```

---

### 3. Get All Rewards

**GET** `/api/reward?page=1&page_size=10`

Retrieve all reward events with pagination.

**Query Parameters:**

- `page` (default: 1)
- `page_size` (default: 10, max: 100)

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "user_name": "John Doe",
      "user_email": "john@example.com",
      "stock_id": 1,
      "stock_symbol": "RELIANCE",
      "stock_name": "Reliance Industries Ltd",
      "quantity": 10.5,
      "stock_price": 2450.75,
      "total_value": 25732.88,
      "event_type": "REWARD",
      "status": "COMPLETED",
      "description": "Performance bonus Q4",
      "created_at": "2025-12-19T10:30:00Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 45,
  "total_pages": 5
}
```

---

### 4. Get Rewards by User

**GET** `/api/reward/user/:userId?page=1&page_size=10`

Retrieve all rewards for a specific user.

**Response:** Same structure as Get All Rewards

---

## User Endpoints

### 1. Get All Users

**GET** `/api/users?page=1&page_size=10`

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "email": "john@example.com",
      "name": "John Doe",
      "phone": "+919876543210",
      "is_active": true,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 100,
  "total_pages": 10
}
```

---

### 2. Get User by ID

**GET** `/api/users/:id`

**Response:** `200 OK`

```json
{
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe",
    "phone": "+919876543210",
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 3. Get Today's Stock Rewards

**GET** `/api/today-stocks/:userId?page=1&page_size=10`

Get all stock rewards issued to a user today.

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "stock_symbol": "RELIANCE",
      "stock_name": "Reliance Industries Ltd",
      "quantity": 10.5,
      "stock_price": 2450.75,
      "total_value": 25732.88,
      "description": "Performance bonus Q4",
      "created_at": "2025-12-19T10:30:00Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 3,
  "total_pages": 1
}
```

---

### 4. Get Historical INR Values

**GET** `/api/historical-inr/:userId?page=1&page_size=10`

Get daily aggregated INR value of rewards for past days.

**Response:** `200 OK`

```json
{
  "data": [
    {
      "date": "2025-12-18",
      "total_value": 45000.5,
      "reward_count": 5
    },
    {
      "date": "2025-12-17",
      "total_value": 32000.0,
      "reward_count": 3
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 30,
  "total_pages": 3
}
```

---

### 5. Get User Statistics

**GET** `/api/stats/:userId`

Get today's rewards summary and current portfolio value.

**Response:** `200 OK`

```json
{
  "today_rewards": [
    {
      "stock_symbol": "RELIANCE",
      "stock_name": "Reliance Industries Ltd",
      "total_shares": 15.5
    },
    {
      "stock_symbol": "TCS",
      "stock_name": "Tata Consultancy Services",
      "total_shares": 8.0
    }
  ],
  "current_portfolio_value": 125430.5
}
```

---

### 6. Get User Portfolio

**GET** `/api/portfolio/:userId?page=1&page_size=10`

Get user's current stock holdings with profit/loss calculations.

**Response:** `200 OK`

```json
{
  "data": [
    {
      "stock_symbol": "RELIANCE",
      "stock_name": "Reliance Industries Ltd",
      "total_quantity": 25.5,
      "average_price": 2400.0,
      "current_price": 2450.75,
      "current_value": 62494.13,
      "investment_cost": 61200.0,
      "profit_loss": 1294.13
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 5,
  "total_pages": 1,
  "total_portfolio_value": 125430.5
}
```

**Note:** All INR values are rounded to 2 decimal places to prevent rounding errors.

---

## Stock Endpoints

### 1. Get All Stocks

**GET** `/api/stocks`

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "symbol": "RELIANCE",
      "name": "Reliance Industries Ltd",
      "exchange": "NSE",
      "current_price": 2450.75,
      "is_active": true,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-12-19T00:00:00Z"
    }
  ]
}
```

---

### 2. Get Stock by ID

**GET** `/api/stocks/:id`

**Response:** `200 OK`

```json
{
  "data": {
    "id": 1,
    "symbol": "RELIANCE",
    "name": "Reliance Industries Ltd",
    "exchange": "NSE",
    "current_price": 2450.75,
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-12-19T00:00:00Z"
  }
}
```

---

### 3. Get Stock by Symbol

**GET** `/api/stocks/symbol/:symbol`

**Response:** Same as Get Stock by ID

---

### 4. Create Stock

**POST** `/api/stocks`

**Request Body:**

```json
{
  "symbol": "RELIANCE",
  "name": "Reliance Industries Ltd",
  "exchange": "NSE",
  "current_price": 2450.75
}
```

**Response:** `201 Created`

```json
{
  "data": {
    "id": 1,
    "symbol": "RELIANCE",
    "name": "Reliance Industries Ltd",
    "exchange": "NSE",
    "current_price": 2450.75,
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 5. Update Stock

**PUT** `/api/stocks/:id`

**Request Body:**

```json
{
  "name": "Reliance Industries Limited",
  "exchange": "NSE",
  "current_price": 2500.0,
  "is_active": true
}
```

**Response:** `200 OK` (same structure as Create Stock)

---

### 6. Delete Stock

**DELETE** `/api/stocks/:id`

Soft delete (sets is_active to false).

**Response:** `200 OK`

```json
{
  "message": "Stock deleted successfully"
}
```

---

## Corporate Action Endpoints

### 1. Create Corporate Action

**POST** `/api/corporate-action`

Create a corporate action (stock split, merger, or delisting).

**Request Body - Stock Split:**

```json
{
  "stock_symbol": "RELIANCE",
  "action_type": "STOCK_SPLIT",
  "split_ratio": 2.0,
  "effective_date": "2025-12-25",
  "description": "1:2 stock split"
}
```

**Request Body - Merger:**

```json
{
  "stock_symbol": "COMPANY_A",
  "action_type": "MERGER",
  "merger_to_symbol": "COMPANY_B",
  "merger_ratio": 0.5,
  "effective_date": "2025-12-25",
  "description": "Merger with Company B at 1:0.5 ratio"
}
```

**Request Body - Delisting:**

```json
{
  "stock_symbol": "OLD_COMPANY",
  "action_type": "DELISTING",
  "effective_date": "2025-12-25",
  "description": "Company delisting from exchange"
}
```

**Response:** `201 Created`

```json
{
  "message": "Corporate action created successfully",
  "data": {
    "id": 1,
    "stock_symbol": "RELIANCE",
    "action_type": "STOCK_SPLIT",
    "split_ratio": 2.0,
    "merger_to_symbol": "",
    "merger_ratio": 0.0,
    "effective_date": "2025-12-25",
    "status": "PENDING",
    "description": "1:2 stock split",
    "affected_users": 150,
    "created_at": "2025-12-19T10:00:00Z",
    "processed_at": null
  }
}
```

---

### 2. Process Corporate Action

**POST** `/api/corporate-action/:id/process`

Execute a pending corporate action.

**Response:** `200 OK`

```json
{
  "message": "Corporate action processed successfully"
}
```

**Effects by Type:**

**Stock Split:**

- Multiplies all user holdings by split_ratio
- Divides stock price by split_ratio
- Example: 1:2 split → 10 shares @ ₹1000 → 20 shares @ ₹500

**Merger:**

- Transfers holdings to target stock with conversion ratio
- Zeros out holdings in source stock
- Deactivates source stock
- Example: 1:0.5 ratio → 10 shares of A → 5 shares of B

**Delisting:**

- Zeros out all user holdings
- Deactivates the stock
- No new rewards can be issued

**Error Responses:**

```json
{
  "error": "corporate action not found"
}
```

```json
{
  "error": "corporate action already processed"
}
```

---

### 3. Get All Corporate Actions

**GET** `/api/corporate-action?page=1&page_size=10`

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "stock_symbol": "RELIANCE",
      "action_type": "STOCK_SPLIT",
      "split_ratio": 2.0,
      "merger_to_symbol": "",
      "merger_ratio": 0.0,
      "effective_date": "2025-12-25",
      "status": "COMPLETED",
      "description": "1:2 stock split",
      "affected_users": 150,
      "created_at": "2025-12-19T10:00:00Z",
      "processed_at": "2025-12-25T00:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 25,
  "total_pages": 3
}
```

---

## Ledger Endpoints

### 1. Get User Ledger Entries

**GET** `/api/ledger/user/:userId`

Get all double-entry ledger records for a user.

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "reward_event_id": 1,
      "stock_id": 1,
      "stock_symbol": "RELIANCE",
      "entry_type": "DEBIT",
      "account_type": "STOCK_UNITS",
      "quantity": 10.5,
      "amount": null,
      "description": "Stock reward credited",
      "created_at": "2025-12-19T10:30:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "reward_event_id": 1,
      "stock_id": null,
      "stock_symbol": "",
      "entry_type": "CREDIT",
      "account_type": "INR_CASH",
      "quantity": null,
      "amount": 25732.88,
      "description": "Cash outflow for stock purchase",
      "created_at": "2025-12-19T10:30:00Z"
    }
  ]
}
```

---

### 2. Get User Stock Holdings

**GET** `/api/ledger/holdings/:userId`

Get current stock holdings with profit/loss.

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "user_name": "John Doe",
      "user_email": "john@example.com",
      "stock_id": 1,
      "stock_symbol": "RELIANCE",
      "stock_name": "Reliance Industries Ltd",
      "total_quantity": 25.5,
      "average_price": 2400.0,
      "current_price": 2450.75,
      "current_value": 62494.13,
      "profit_loss": 1294.13,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-12-19T10:30:00Z"
    }
  ]
}
```

---

### 3. Get All Stock Holdings

**GET** `/api/ledger/holdings`

Get all users' stock holdings.

**Response:** Same structure as Get User Stock Holdings

---

### 4. Get Account Summary

**GET** `/api/ledger/summary?userId=1`

Get summary of all account types for a user (optional userId filter).

**Response:** `200 OK`

```json
{
  "data": [
    {
      "user_id": 1,
      "user_name": "John Doe",
      "account_type": "STOCK_UNITS",
      "total_quantity": 35.5,
      "total_amount": null
    },
    {
      "user_id": 1,
      "user_name": "John Doe",
      "account_type": "INR_CASH",
      "total_quantity": null,
      "total_amount": 125430.5
    }
  ]
}
```

---

## Common Response Codes

- `200 OK` - Successful GET/PUT/DELETE request
- `201 Created` - Successful POST request
- `400 Bad Request` - Invalid request parameters or validation failure
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Pagination

All list endpoints support pagination:

**Query Parameters:**

- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)

**Response includes:**

- `data` - Array of items
- `page` - Current page
- `page_size` - Items per page
- `total_count` - Total items
- `total_pages` - Total pages

---

## Data Precision

- **Stock Quantities:** 6 decimal places (NUMERIC 18,6)
- **INR Amounts:** 2 decimal places (rounded from NUMERIC 18,4)
- **Stock Prices:** 4 decimal places (NUMERIC 18,4)
- **Percentages:** 4 decimal places (NUMERIC 5,4)

All INR values in API responses are rounded to 2 decimal places to prevent floating-point rounding errors.

# Stocky Backend Assignment

A robust stock reward management system built with Go, featuring double-entry ledger accounting, corporate actions handling, and comprehensive portfolio tracking.

## 🎯 Project Overview

This backend system manages stock rewards for users with enterprise-grade features including:

- **Double-entry ledger** accounting for all transactions
- **Stock reward issuance** with duplicate detection and idempotency
- **Portfolio management** with real-time profit/loss calculations
- **Corporate actions** (stock splits, mergers, delisting)
- **Reward adjustments/refunds** with full audit trail
- **Comprehensive reporting** with pagination support

## ✨ Features

### Core Functionality

- ✅ Issue stock rewards to users with automatic ledger entries
- ✅ Track user portfolios with current values and profit/loss
- ✅ Handle corporate actions (stock splits, mergers, delistings)
- ✅ Refund/adjust previously issued rewards
- ✅ Double-entry bookkeeping for financial accuracy
- ✅ Prevent duplicate rewards (time-based + idempotency keys)

### Technical Features

- ✅ RESTful API with pagination
- ✅ PostgreSQL with ACID transactions
- ✅ Comprehensive error handling
- ✅ Structured logging
- ✅ Database migrations
- ✅ Hot-reload development environment
- ✅ INR rounding precision (2 decimal places)

## 🛠️ Technologies Used

- **Language:** Go 1.23.4
- **Web Framework:** Gin
- **Database:** PostgreSQL 14+
- **Logging:** Logrus
- **Configuration:** godotenv
- **Development:** Air (hot-reload)

## 📋 Prerequisites

Before running this project, ensure you have:

- **Go** 1.23.4 or higher ([Download](https://golang.org/dl/))
- **PostgreSQL** 14 or higher ([Download](https://www.postgresql.org/download/))
- **Git** ([Download](https://git-scm.com/downloads))

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/jayzadafiya/Stocky-Backend-Assignment.git
cd Stocky-Backend-Assignment
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up Database

#### Create PostgreSQL Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE stocky_db;

# Create user (optional)
CREATE USER stocky_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE stocky_db TO stocky_user;

# Exit
\q
```

### 4. Configure Environment Variables

Create a `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=stocky_db
DB_SSLMODE=disable

# Server Configuration
PORT=8080
GIN_MODE=release

# Logging
LOG_LEVEL=info
```

### 5. Run Database Migrations

```bash
# Using the migration tool
go run cmd/migrate/main.go

# Or using make
make migrate
```

This will create all required tables:

- `users`
- `stocks`
- `reward_events`
- `ledger_entries`
- `user_stock_holdings`
- `fee_configurations`
- `corporate_actions`

### 6. Seed Initial Data (Optional)

```bash
# Load sample users and stocks from data/ folder
# (You can manually import or create a seed script)
```

### 7. Run the Application

#### Production Mode

```bash
# Using Go
go run main.go

# Or using make
make run
```

#### Development Mode (with hot-reload)

```bash
# Install Air (first time only)
go install github.com/air-verse/air@latest

# Run with hot-reload
make watch

# Or using PowerShell script
.\run.ps1 watch
```

The server will start at `http://localhost:8080`

### 8. Verify Installation

```bash
# Health check
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok",
  "message": "Stocky Backend Assignment API is running"
}
```

## 📁 Project Structure

```
stocky-backend/
├── cmd/
│   └── migrate/           # Database migration tool
├── config/
│   ├── database.go        # Database connection
│   └── logger.go          # Logging configuration
├── data/
│   ├── stocks.json        # Sample stock data
│   └── users.json         # Sample user data
├── features/
│   ├── corporate_action/  # Corporate actions module
│   │   ├── corporate_action.model.go
│   │   ├── corporate_action.service.go
│   │   ├── corporate_action.handler.go
│   │   └── corporate_action.routes.go
│   ├── ledger/           # Ledger entries module
│   ├── reward/           # Reward management module
│   ├── stock/            # Stock management module
│   └── user/             # User management module
├── migrations/           # SQL migration files
│   ├── 001_create_users_table.sql
│   ├── 002_create_stocks_table.sql
│   ├── 003_create_reward_events_table.sql
│   ├── 004_create_ledger_entries_table.sql
│   ├── 005_create_user_stock_holdings_table.sql
│   ├── 006_create_fee_configurations_table.sql
│   └── 007_create_corporate_actions_table.sql
├── .air.toml            # Hot-reload configuration
├── .env.example         # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
├── main.go              # Application entry point
├── Makefile             # Build commands
├── run.ps1              # PowerShell run script
├── API_DOCUMENTATION.md      # Complete API docs
├── DATABASE_SCHEMA.md        # Database schema docs
├── ARCHITECTURE.md           # System architecture docs
└── README.md            # This file
```

## 🔌 API Endpoints

### Base URL: `http://localhost:8080/api`

### Quick Reference

| Module                | Method | Endpoint                        | Description             |
| --------------------- | ------ | ------------------------------- | ----------------------- |
| **Health**            | GET    | `/health`                       | Health check            |
| **Rewards**           | POST   | `/reward`                       | Create stock reward     |
|                       | POST   | `/reward/adjust`                | Refund/adjust reward    |
|                       | GET    | `/reward`                       | List all rewards        |
|                       | GET    | `/reward/user/:userId`          | Get user rewards        |
| **Users**             | GET    | `/users`                        | List all users          |
|                       | GET    | `/users/:id`                    | Get user by ID          |
|                       | GET    | `/today-stocks/:userId`         | Today's rewards         |
|                       | GET    | `/historical-inr/:userId`       | Historical INR values   |
|                       | GET    | `/stats/:userId`                | User statistics         |
|                       | GET    | `/portfolio/:userId`            | User portfolio          |
| **Stocks**            | GET    | `/stocks`                       | List all stocks         |
|                       | GET    | `/stocks/:id`                   | Get stock by ID         |
|                       | GET    | `/stocks/symbol/:symbol`        | Get stock by symbol     |
|                       | POST   | `/stocks`                       | Create stock            |
|                       | PUT    | `/stocks/:id`                   | Update stock            |
|                       | DELETE | `/stocks/:id`                   | Delete stock            |
| **Corporate Actions** | POST   | `/corporate-action`             | Create corporate action |
|                       | POST   | `/corporate-action/:id/process` | Process action          |
|                       | GET    | `/corporate-action`             | List all actions        |
| **Ledger**            | GET    | `/ledger/user/:userId`          | User ledger entries     |
|                       | GET    | `/ledger/holdings/:userId`      | User stock holdings     |
|                       | GET    | `/ledger/holdings`              | All holdings            |
|                       | GET    | `/ledger/summary`               | Account summary         |

📖 **For detailed API documentation with request/response examples, see [API_DOCUMENTATION.md](API_DOCUMENTATION.md)**

## 💻 Development Commands

### Using Makefile

```bash
# Run the application
make run

# Run with hot-reload
make watch

# Run database migrations
make migrate

# Build the application
make build

# Clean build artifacts
make clean
```

### Using PowerShell Script (Windows)

```powershell
# Development mode
.\run.ps1 dev

# Production mode
.\run.ps1 prod

# Hot-reload mode
.\run.ps1 watch

# Database migration
.\run.ps1 migrate
```

## 🧪 Testing

### Manual Testing with curl

```bash
# Create a reward
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "stock_symbol": "RELIANCE",
    "quantity": 10.5,
    "description": "Performance bonus"
  }'

# Get user portfolio
curl http://localhost:8080/api/portfolio/1

# Process corporate action
curl -X POST http://localhost:8080/api/corporate-action/1/process
```

### Using Postman

Import the Postman collection from:
🔗 [Postman Workspace](https://www.postman.com/grey-resonance-948392/workspace/stocky-backend-assignment)

## 📊 Database Schema

The system uses PostgreSQL with the following tables:

- **users** - User accounts
- **stocks** - Stock/security information
- **reward_events** - All reward transactions
- **ledger_entries** - Double-entry accounting ledger
- **user_stock_holdings** - Current user holdings
- **corporate_actions** - Stock splits, mergers, delistings
- **fee_configurations** - Transaction fees

📖 **For complete schema documentation, see [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md)**

## 🏗️ Architecture

The system is built with a **feature-based architecture**:

- **Modular design** - Each feature is self-contained
- **Double-entry ledger** - Financial accuracy and audit trail
- **ACID transactions** - Data consistency guaranteed
- **Scalable** - Ready for horizontal scaling
- **Edge case handling** - Duplicate prevention, concurrent access, rounding errors

📖 **For architecture details and scaling considerations, see [ARCHITECTURE.md](ARCHITECTURE.md)**

## 🔐 Security Features

- ✅ SQL injection prevention (parameterized queries)
- ✅ Input validation with Gin bindings
- ✅ Database constraints and checks
- ✅ Transaction atomicity
- ✅ Idempotency key support

## 🐛 Common Issues & Solutions

### Issue: Database connection failed

**Solution:** Check your PostgreSQL is running and `.env` credentials are correct

```bash
# Check PostgreSQL status
pg_ctl status

# Restart PostgreSQL
pg_ctl restart
```

### Issue: Port 8080 already in use

**Solution:** Change the PORT in `.env` file or stop the conflicting process

```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:8080 | xargs kill -9
```

### Issue: Migration failed

**Solution:** Drop and recreate the database

```sql
DROP DATABASE IF EXISTS stocky_db;
CREATE DATABASE stocky_db;
```

Then re-run migrations.

## 📝 Environment Variables Reference

| Variable      | Description       | Default     | Required |
| ------------- | ----------------- | ----------- | -------- |
| `DB_HOST`     | PostgreSQL host   | `localhost` | Yes      |
| `DB_PORT`     | PostgreSQL port   | `5432`      | Yes      |
| `DB_USER`     | Database user     | `postgres`  | Yes      |
| `DB_PASSWORD` | Database password | -           | Yes      |
| `DB_NAME`     | Database name     | `stocky_db` | Yes      |
| `DB_SSLMODE`  | SSL mode          | `disable`   | Yes      |
| `PORT`        | Server port       | `8080`      | No       |
| `GIN_MODE`    | Gin mode          | `release`   | No       |
| `LOG_LEVEL`   | Log level         | `info`      | No       |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

**Jay Zadafiya**

- GitHub: [@jayzadafiya](https://github.com/jayzadafiya)
- Postman: [Workspace](https://www.postman.com/grey-resonance-948392/workspace/stocky-backend-assignment)

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Web framework: [Gin](https://gin-gonic.com/)
- Database: [PostgreSQL](https://www.postgresql.org/)
- Logging: [Logrus](https://github.com/sirupsen/logrus)

---

## 📚 Additional Documentation

- 📖 [API Documentation](API_DOCUMENTATION.md) - Complete API specifications
- 📖 [Database Schema](DATABASE_SCHEMA.md) - Database design and relationships
- 📖 [Architecture](ARCHITECTURE.md) - System design and scaling guide
- 📖 [Quickstart Guide](QUICKSTART.md) - Quick reference guide

---

**Made with ❤️ for efficient stock reward management**

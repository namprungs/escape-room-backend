## Setup Instructions

#### 1. Clone the Repository
```bash
git clone https://github.com/FieldPs/escape-room-backend
cd escape-room-backend
```

#### 2. Configure Environment Variables
Create a `.env` file in the project root with the following content
```
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=puzzle_db
DB_PORT=5432
DB_SSLMODE=disable

# PostgreSQL Configuration
POSTGRES_PASSWORD=secret
POSTGRES_USER=postgres
POSTGRES_DB=puzzle_db

# JWT Configuration
JWT_SECRET=ThisIsASecretKeyForJWT

# Application Configuration
APP_PORT=8080
APP_ENV=development
```

#### 3. Initial setup
1. Start PostgreSQL `docker-compose up -d`
2. Install Go Dependencies `go mod download`

#### 4. run project locally
`go run cmd/main.go`

## API Endpoints

| Method | Endpoint             | Description                  | Auth | Payload Example                     |
|--------|----------------------|------------------------------|------|-------------------------------------|
| GET    | `/healthz`           | Check app health             |No    | None                               |
| POST   | `/api/v1/register`   | Register a new user          |No    | `{"username": "test", "password": "pass123"}` |
| POST   | `/api/v1/login`      | Login and get JWT            |No    | `{"username": "test", "password": "pass123"}` |
| GET    | `/api/v1/stats`      | Get Stats of User            | ✅  | None                                 |
| POST   | `/api/v1/submit_answer`| send puzzle answer         | ✅  | `{"Puzzle_id" : 1, "answer" : "1234"}` |

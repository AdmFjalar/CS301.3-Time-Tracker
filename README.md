# CS301.3-Time-Tracker

## Instructions for Setting Up the Server

1. **Clone the repository:**
   ```sh
   git clone https://github.com/AdmFjalar/CS301.3-Time-Tracker.git
   cd CS301.3-Time-Tracker
   ```

2. **Set up environment variables:**
   Create a `.env` file in the `backend` directory with the following content:
   ```sh
   ADDR=:8080
   EXTERNAL_URL=localhost:8080
   FRONTEND_URL=http://localhost:5173
   DB_ADDR=user:password@tcp(localhost:3306)/thymeflies
   DB_MAX_OPEN_CONNS=30
   DB_MAX_IDLE_CONNS=30
   DB_MAX_IDLE_TIME=15m
   DB_USER=user
   DB_PASSWORD=password
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=thymeflies
   REDIS_ADDR=localhost:3306
   REDIS_PW=""
   REDIS_DB=0
   REDIS_ENABLED=false
   ENV=development
   FROM_EMAIL=from@email.com
   SENDGRID_API_KEY=yourAPIkey
   AUTH_BASIC_USER=admin
   AUTH_BASIC_PASS=admin
   AUTH_TOKEN_SECRET=tokenexample
   RATELIMITER_REQUESTS_COUNT=2
   RATE_LIMITER_ENABLED=true
   JWT_SECRET=not-so-secret-now-is-it?
   JWT_EXPIRATION_IN_SECONDS=604800
   SERVER_ADDRESS=:8080

   ```

3. **Install dependencies:**
   ```sh
   cd backend
   go mod tidy
   ```

4. **Run the server:**
   ```sh
   go run cmd/api/main.go
   ```

5. **Access the server:**
   The server will be running at `http://localhost:8080`. You can use tools like `curl` or Postman to interact with the API endpoints.

## Running Tests

1. **Run unit tests:**
   ```sh
   go test ./...
   ```

2. **Run integration tests:**
   ```sh
   go test -tags=integration ./...
   ```

## Deploying the Application

1. **Build the application:**
   ```sh
   go build -o bin/app cmd/api/main.go
   ```

2. **Run the built application:**
   ```sh
   ./bin/app
   ```

3. **Deploy to a server:**
   Copy the built application and the `.env` file to your server and run the application as described above.

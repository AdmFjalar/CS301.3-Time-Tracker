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

4. **Create MySQL user:**
   Create a MySQL user and add the corresponding details in the .env file.

5. **Run the SQL instantiation:**
   Run the SQL file in the main folder to create a database and all the tables related to it.

6. **Run the server:**
   ```sh
   go run cmd/api/main.go
   ```

4. **Run the frontend:**
   In a separate terminal, run
      ```sh
   cd frontend
   npm install
   npm start
   ```

7. **Access the server:**
   The website will be running at `http://localhost:3000`.
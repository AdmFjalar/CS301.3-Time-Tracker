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
   DB_USER=root
   DB_PASSWORD=mypassword
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=ecom
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

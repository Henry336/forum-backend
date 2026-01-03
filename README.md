# CVWO Forum Backend (Go)

A backend server built with Golang and PostgreSQL for the CVWO: Gossip with Go assignment.

## Features
- REST API implementation (GET, POST, PATCH, DELETE)
- PostgreSQL database integration
- CORS enabled for frontend communication

## How to Run

1. **Prerequisites:**
   - Go installed
   - PostgreSQL installed and running on port 5432

2. **Database Setup:**
   - Create a database named `cvwo_db`.
   - Run the following SQL to create the table:
     ```sql
     CREATE TABLE topics (
       id SERIAL PRIMARY KEY,
       title TEXT NOT NULL,
       description TEXT NOT NULL
     );
     ```

3. **Start the Server:**
   ```bash
   go run main.go
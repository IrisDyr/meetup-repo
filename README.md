# meetup-repo

## Steps to Run the Setup

### Step 1: Deploy OpenTelemetry, Tempo, and Grafana
1. Navigate to the `otel-tempo-grafana` directory:
   ```bash
   cd otel-tempo-grafana
   ```
2. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

### Step 2: Run the Go Application
1. Navigate to the `otel-go-app` directory:
   ```bash
   cd otel-go-app
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the application:
   ```bash
   go run main.go
   ```

### Step 3: Check the Traces
- Open Grafana in your browser:
  ```
  http://localhost:3000
  ```

### Step 4: Clean Up Your Docker Environment
- Stop and remove the containers:
  ```bash
  docker-compose down
  

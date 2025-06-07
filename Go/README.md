# Load Balancer

A flexible and configurable load balancer written in Go, designed for high availability and performance.

## Features

- Multiple load balancing algorithms
- Health checking with configurable thresholds
- Rate limiting per IP
- Smart mode for automatic failover
- Real-time metrics dashboard
- State history tracking
- Authentication for metrics dashboard
- CORS support
- SQLite database for metrics storage

## Quick Start

### 1. Pull Docker Image
```bash
docker pull kaiz0708/load-balancer:latest
```

### 2. Configuration
Create a `config.yaml` file:
```yaml
# Load balancing algorithm
algorithm: "ROUND_ROBIN"  # Required field

# List of backend servers
backends:
  - url: "http://backend1:8080"
    weight: 1
    healthPath: "/api/health"
  - url: "http://backend2:8080"
    weight: 3
    healthPath: "/api/health"

# Health check settings
consecutiveFails: 10     # Number of consecutive failures before marking backend as unhealthy
consecutiveSuccess: 10   # Number of consecutive successes before marking backend as healthy
failRate: 0.5           # Failure rate threshold (0.0 to 1.0)
timeOutRate: 10         # Timeout break duration in seconds
timeOutDelay: 5         # Timeout delay in seconds

# Authentication for metrics dashboard
auth:
  username: "your-username"     # Required field
  password: "your-password"  # Required field

# Smart mode for automatic failover
smartMode: true

# Rate limiting
rateLimit: 100         # Maximum requests per second per IP
```

### 3. Run with Docker
```bash
docker run -d \
  --name load-balancer \
  -p 8080:8080 \
  -v /path/to/your/config.yaml:/app/config.yaml \
  kaiz070/load-balancer:latest
```

## Default Values

If not configured in `config.yaml`, the following default values will be used:

### Health Check Settings
- `consecutiveFails`: 10 (Consecutive failures before marking backend as unhealthy)
- `consecutiveSuccess`: 10 (Consecutive successes to mark backend as healthy)
- `failRate`: 0.5 (Allowed failure rate, from 0.0 to 1.0)
- `timeOutRate`: 10 (Timeout break duration in seconds)
- `timeOutDelay`: 5 (Timeout delay in seconds)
- `healthCheckInterval`: 1 (Interval between health checks in seconds)

### Rate Limiting
- `rateLimit`: 100 (Maximum requests per second per IP)

### Backend Settings
- `weight`: 1 (Default weight for each backend server)

### Smart Mode
- `smartMode`: false (Smart mode disabled by default)

### Required Fields
- `algorithm`: Must be specified
- `auth.username`: Must be specified
- `auth.password`: Must be specified
- `backends`: At least one backend must be specified

## Load Balancing Algorithms

### Standard Algorithms
- `ROUND_ROBIN`: Distributes requests evenly across servers
- `LEAST_CONNECTION`: Selects server with least active connections
- `WEIGHTED_ROUND_ROBIN`: Round robin with server weights
- `WEIGHTED_LEAST_CONNECTION`: Least connection with server weights
- `RANDOM`: Random server selection
- `WEIGHTED_RANDOM`: Random selection with server weights
- `IP_HASH`: Distributes based on client IP

### Smart Mode Algorithms
- `WEIGHTED_SUCCESS_RATE_BALANCER`: Used when many servers are failing
- `LOW_LATENCY_WEIGHTED_BALANCER`: Used when high latency is detected

## API Endpoints

### Load Balancing
- `GET /`: Main load balancer endpoint
- `GET /change-load-balancer`: Change load balancing algorithm

### Metrics & Monitoring
- `GET /metrics`: Metrics dashboard (requires authentication)
- `POST /login-metrics`: Login to metrics dashboard
- `GET /state-history`: Get error history
- `POST /delete-state-history`: Delete error history
- `POST /reset-metrics`: Reset metrics for a backend

## System States

- `Stable`: System operating normally
- `AllFailed`: All backend servers are unhealthy
- `ManyFailed`: More than half of backend servers are unhealthy
- `HighLatency`: More than half of backend servers have high latency
- `Recovery`: Backend server is recovering
- `Unhealthy`: Backend server is unhealthy

## Error Types

- `connection refused`: Connection refused error
- `no such host`: Host not found error
- `context canceled`: Request was canceled by the client
- `context deadlineExceeded`: Request exceeded the configured timeout

## Environment Variables

- `CONFIG_FILE`: Path to config file ("config.yaml")
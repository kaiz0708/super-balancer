# Load Balancer

A flexible and configurable load balancer written in Go.

## Quick Start

1. Pull the Docker image:
```bash
docker pull your-username/load-balancer:latest
```

2. Create your configuration file (config.yaml):
```yaml
# Load balancing algorithm to use
# Available options: ROUND_ROBIN, LEAST_CONNECTION, WEIGHTED_LEAST_CONNECTION,
#                   WEIGHTED_ROUND_ROBIN, RANDOM, WEIGHTED_RANDOM, IP_HASH
algorithm: "ROUND_ROBIN"

# List of backend servers
backends:
  - url: "http://your-backend1:8080"
    weight: 1
    healthPath: "/api/health"
  - url: "http://your-backend2:8080"
    weight: 3
    healthPath: "/api/health"

# Health check settings
consecutiveFails: 100     # Number of consecutive failures before marking backend as unhealthy
consecutiveSuccess: 100   # Number of consecutive successes before marking backend as healthy
failRate: 0.5          # Failure rate threshold (0.0 to 1.0)
timeOutBreak: 100       # Timeout break duration in seconds
timeOutDelay: 5        # Timeout delay in seconds

# Authentication for metrics dashboard
auth:
  username: "admin"     # Your desired username
  password: "your-password"  # Your desired password

# Enable/disable smart mode for automatic failover
smartMode: true

# Rate limiting
rateLimit: 100         # Maximum requests per second per IP
```

3. Run the load balancer:
```bash
docker run -d \
  --name load-balancer \
  -p 8080:8080 \
  -v /path/to/your/config.yaml:/app/config.yaml \
  your-username/load-balancer:latest
```

## Features

- Multiple load balancing algorithms
- Health checking
- Rate limiting
- Smart mode for automatic failover
- Metrics dashboard
- Error history tracking
- Authentication for metrics dashboard

## Metrics Dashboard

Access the metrics dashboard at `http://localhost:8080/metrics`

## API Endpoints

- `GET /`: Main load balancer endpoint
- `GET /metrics`: Metrics dashboard
- `POST /login-metrics`: Login to metrics dashboard
- `GET /change-load-balancer`: Change load balancing algorithm
- `GET /error-history`: Get error history
- `POST /delete-error-history`: Delete error history
- `POST /reset-metrics`: Reset metrics for a backend 
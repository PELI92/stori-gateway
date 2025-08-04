# 🧱 Build a Minimal API Gateway

## 🎯 Objective

Design and implement a lightweight API Gateway / Reverse Proxy in **Go** that:

- Routes incoming HTTP requests to downstream services based on path.
- Performs basic validation (`X-Request-ID`).
- Handles simple authentication via `x-api-key`.
- Supports live configuration reloads.


---

## 🛠️ Tech Stack

| Component | Purpose |
|----------|---------|
| **Gin** | HTTP server and routing — fast, minimal, and familiar. |
| **ReverseProxy** (`net/http/httputil.ReverseProxy`) | Standard Go way to proxy requests, header/body-safe. |
| **Koanf** | Config management with hot reload support. |

---

## 📂 Directory Structure (simplified)

```
stori-gateway/
├── cmd/
│   └── gateway/              # Entry point (main.go)
│       └── main.go
│
├── config/
│   └── config.yaml           # Runtime configuration (services, API key, etc.)
│
├── internal/                 # Application core
│   ├── config/               # Configuration loader and provider
│   ├── middleware/           # Gin middlewares (auth, validation, logging)
│   └── proxy/                # Reverse proxy logic and request forwarding
│
├── docs/
│   └── swagger.go            # Swagger spec definition
│
├── mock/                     # Mock backend service for testing
│
├── test/
│   └── load-test.yml         # Artillery load test scenario
│
├── docker-compose.yml        # Orchestration of all services
├── Dockerfile                # Gateway Dockerfile
├── go.mod                    # Project dependencies
└── README.md                 # Project documentation

```

## 🚀 How to Run

```bash
  docker compose up --build
```
Notes: 

This will create a container for this service plus 2 mock services to test it: user and auth.

When running locally, hot reload can be tested by changing config/condig.yaml for simplicity. In a cloud environment this should be configured for a VM volume file, S3 file, secret or via DB

##  Examples

#### ✅  OK
```bash
  curl -i -H "X-Request-ID: abc123" -H "x-api-key: supersecretkey" http://localhost:8080/api/auth/hello
```

#### 🚫 Missing X-Request-ID
```bash
  curl -i -H "x-api-key: supersecretkey" http://localhost:8080/api/auth/hello
```

#### 🚫 Missing x-api-key
```bash
  curl -i -H "X-Request-ID: test" http://localhost:8080/api/auth/hello
```

#### 🚫 Incorrect API key
```bash
  curl -i -H "X-Request-ID: test" -H "x-api-key: wrong" http://localhost:8080/api/auth/hello
```

#### ✅ Swagger json (can be imported to postman for testing)
```bash
  curl -i http://localhost:8080/swagger.json
```

#### 🚀 Load testing

When the app is running cd to test directory:

```bash
    cd test
```

and execute with artillery:

```bash
  artillery run load-test.yml
```

To install artillery:
```bash
  npm install -g artillery
```

This will test the app with GET and POST (w/JSON) at 1000 req/s each, for 10 seconds. As configured, it can easily handler 2000 req/s combined. For a more demanding context I would suggest adding multiple instances and a load balancer.

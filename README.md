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
.
├── cmd/gateway              # Application entrypoint
├── config/config.yaml       # Route + API key config
├── internal/
│   ├── config               # Config loading + hot reload
│   ├── middleware           # Validation + logging
│   └── proxy                # Reverse proxy logic
├── Dockerfile
└── docker-compose.yml
```

## 🚀 How to Run

```bash
  docker compose up --build
```
Notes: 

This will create a container for this service plus 2 mock services to test it: user and auth.

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

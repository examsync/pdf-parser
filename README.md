# PDF Parser Service (pdf-parser)

A structured, modular Go web microservice following Clean and Layered Architecture principles. This project mirrors the architecture of `exam-mgmt-service` and provides a robust, decoupled structure for managing parsed PDF metadata.

---

## 🚀 Core Technologies & Dependencies

- **Web Framework**: [Echo v5 (v5.2.1)](https://github.com/labstack/echo) — Lightweight, fast Go web framework for HTTP routing and context handling.
- **ORM / Database Client**: [GORM (v1.31.2)](https://gorm.io/) with [PostgreSQL Driver (v1.6.0)](https://github.com/go-gorm/postgres) — Object Relational Mapper for PostgreSQL schema auto-migrations and persistence.
- **Configuration Management**: [Viper (v1.21.0)](https://github.com/spf13/viper) — Environment-aware and flexible configuration management.
- **Structured Logging**: [Logrus (v1.9.4)](https://github.com/sirupsen/logrus) — High-performance JSON-formatted structured logging.

---

## 📁 Directory Structure

```text
├── cmd/
│   └── server/
│       ├── handler.go              # MVC dependency wiring and Echo HTTP routing setup
│       ├── main.go                 # App entry point, logger boot, config, and database init
│       └── server.go               # net/http server configuration and graceful shutdown
├── config/
│   └── app.yaml                    # Application config file (server port and database settings)
├── internal/
│   ├── controllers/                # HTTP presentation layer (processes inputs, returns JSON payloads)
│   ├── DTOS/                       # Data Transfer Objects for serialization/deserialization
│   ├── models/                     # Business domain entities (GORM database models)
│   ├── repositories/               # Repository interfaces and storage implementations (GORM queries)
│   └── services/                   # Business logic and service orchestrators
├── utils/
│   ├── common/                     # Shared cross-cutting utility functions
│   ├── config/                     # Configuration structure parser and validation logic
│   ├── database/                   # GORM connection setup, ping validation, and close wrappers
│   └── logger/                     # Logrus JSON logger wrappers and initialization
├── go.mod                          # Go modules dependency descriptor
└── go.sum                          # Go modules checksum validation
```

Detailed file mapping:
- **Cmd Layer**:
  - [cmd/server/main.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/cmd/server/main.go)
  - [cmd/server/handler.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/cmd/server/handler.go)
  - [cmd/server/server.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/cmd/server/server.go)
- **Config**:
  - [config/app.yaml](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/config/app.yaml)
- **Internal Layer (Business & Domain Logic)**:
  - [internal/controllers/controller.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/internal/controllers/controller.go)
  - [internal/DTOS/dto.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/internal/DTOS/dto.go)
  - [internal/models/models.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/internal/models/models.go)
  - [internal/repositories/repo.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/internal/repositories/repo.go)
  - [internal/services/services.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/internal/services/services.go)
- **Utils**:
  - [utils/common/common.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/utils/common/common.go)
  - [utils/config/config.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/utils/config/config.go)
  - [utils/database/database.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/utils/database/database.go)
  - [utils/logger/logger.go](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/utils/logger/logger.go)

---

## 🏗️ Architectural Layers

The service enforces a clean separation of concerns and uses a **concrete dependency injection pattern**:

1. **Presentation Layer (`internal/controllers/` / `cmd/server/handler.go`)**:
   - Decoupled Echo controllers processing HTTP inputs, binding request variables, and calling services.
2. **DTO Layer (`internal/DTOS/`)**:
   - Serialization templates that define payloads for request/response bodies.
3. **Domain Layer (`internal/models/`)**:
   - Domain structures representing core entities mapped to database schemas (e.g. `ExamNotification` database entity).
4. **Service Layer (`internal/services/`)**:
   - Implements core business logic workflows and coordination rules.
5. **Infrastructure & Storage (`internal/repositories/` / `utils/database/`)**:
   - Database connection management, auto-migrations, and data access query execution.

---

## 🔌 API Endpoints

Once the application is running, the following endpoints are available:

* **Get All Exam Notifications (`GET /notifications`)**:
  - Fetches the collection of exam notifications from the PostgreSQL database.
  - Automatically runs GORM auto-migrations and seeds the default BPSSC Company Commander 2026 notification if the table is empty.
* **Health Check (`GET /health`)**:
  - Returns `{"status": "healthy"}` for container/service health diagnostics.

---

## ⚙️ Setup & Configuration

### Prerequisites
- **Go**: Version `1.25.6` or higher.
- **PostgreSQL**: An active database instance.

### Running Locally

1. Download dependencies:
   ```bash
   go mod download
   ```

2. Edit connection details in [config/app.yaml](file:///Users/sharma/go/src/github.com/examsync/pdf-parser/config/app.yaml):
   ```yaml
   server:
     port: 8080

   database:
     host: localhost
     port: 5432
     user: postgres
     password: postgres
     dbname: examsync
     sslmode: disable
   ```

3. Run the application:
   ```bash
   go run cmd/server/*.go
   ```

---

## 🛑 Graceful Shutdown Flow

The server implements signal-based shutdown listening to SIGINT/SIGTERM interrupts:
1. Stops the HTTP listener immediately, rejecting new connection handshakes.
2. Gracefully waits up to **10 seconds** for any running HTTP requests to complete.
3. Closes the GORM PostgreSQL connection pool to finalize any active transactions cleanly and release system resources.
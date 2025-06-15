# Loan Service (LOS)

A comprehensive loan management system that handles the entire loan lifecycle from proposal to disbursement.

## Overview

Loan Service adalah aplikasi backend yang menyediakan API untuk mengelola pinjaman (loans), mulai dari pengajuan, persetujuan, pendanaan, hingga pencairan. Sistem ini dirancang dengan arsitektur yang bersih dan modular untuk memudahkan pengembangan dan pemeliharaan.

## Features

- **Loan Creation**: Create new loan proposals with borrower information, principal amount, rate, and ROI
- **Loan Approval**: Validate and approve loan proposals
- **Investment Management**: Add investments to approved loans
- **Agreement Generation**: Generate loan agreement letters
- **Loan Disbursement**: Disburse fully funded loans
- **Loan Querying**: Retrieve loans by ID, borrower, or state

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Steps

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/los-technical.git
   cd los-technical
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Build the application:
   ```
   go build -o loan-service ./cmd/api
   ```

4. Run the application:
   ```
   ./loan-service
   ```

## API Documentation

Berikut adalah daftar endpoint API yang tersedia dalam aplikasi ini:

### Loan Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/loans` | Create a new loan proposal |
| GET | `/loans/:id` | Get loan by ID |
| POST | `/loans/:id/approve` | Approve a loan |
| POST | `/loans/:id/invest` | Add investment to a loan |
| POST | `/loans/:id/disburse` | Disburse a loan |
| POST | `/loans/:id/agreement` | Generate agreement letter |
| GET | `/loans/borrower/:borrowerId` | Get loans by borrower |
| GET | `/loans/state/:state` | Get loans by state |

### Request/Response Examples

#### Create Loan

Request:
```
POST /loans

{
  "borrowerId": "user123",
  "principalAmount": 10000,
  "rate": 5.5,
  "roi": 10
}
```

Response:
```
{
  "message": "Loan created successfully",
  "data": {
    "id": "loan123",
    "borrowerId": "user123",
    "principalAmount": 10000,
    "rate": 5.5,
    "roi": 10,
    "state": "PROPOSED",
    "createdAt": "2023-01-01T12:00:00Z",
    "updatedAt": "2023-01-01T12:00:00Z"
  }
}
```

## Project Structure

Struktur proyek mengikuti prinsip Clean Architecture dengan pemisahan yang jelas antara domain, use case, dan infrastruktur:

```
.
├── cmd/
│   └── api/                  # Application entry point
│       └── main.go
├── internal/
│   ├── api/                  # API layer
│   │   └── handler/          # HTTP handlers
│   │       └── loan/
│   ├── domain/               # Domain models and interfaces
│   │   ├── loan/
│   │   └── response/
│   ├── infrastructure/       # External implementations
│   │   ├── email/
│   │   └── repository/
│   ├── pkg/                  # Shared utilities
│   │   └── utils/
│   └── usecase/              # Business logic
│       └── loan/
├── go.mod
├── go.sum
└── README.md
```

## Loan Lifecycle

Aplikasi ini mengelola siklus hidup pinjaman melalui beberapa state:

1. **PROPOSED**: Initial state when a loan is created
2. **APPROVED**: Loan has been validated and approved
3. **INVESTED**: Loan has been fully funded by investors
4. **DISBURSED**: Loan has been disbursed to the borrower

Each state transition requires specific validations and actions as implemented in the service layer.

## Development

### Running Tests

To run the tests:

```
go test ./...
or
go test ./... -cover -v -covermode=count -coverprofile=coverage.out 2>&1
```

### Generating Swagger Documentation

```bash
swag init -g cmd/api/main.go -o docs
```

### Adding New Features

When adding new features:

1. Define domain models and interfaces
2. Implement business logic in use cases
3. Create API handlers
4. Update repository implementations as needed

# PYRE PROMOTION SERVICE

This repository hosts a **promotion service** built with Go, featuring a microservices architecture that includes:

- gRPC communication between micro-services
- HTTP communication for front-end
- Kafka for message queues
- PostgreSQL for data storage
- Redis for caching
- OpenTelemetry for distributed tracing
- Health monitoring and metrics

---

## 📚 Table of Contents

- [Overview](#-overview)
- [Architecture](#-architecture)
- [Getting Started](#-getting-started)
- [Deployment](#-deployment)
- [Usage](#-usage)

---

## 🧭 Overview

The Pyre Promotion Service is a robust microservice designed to handle promotional campaigns and discounts in an e-commerce system. It provides:

- Promotion management and validation
- Discount calculation and application
- Event-driven architecture using Kafka
- High-performance caching with Redis
- Distributed tracing with OpenTelemetry
- Better query performance with concurrency

---

## 🏗 Architecture

The service is built with the following key components:

- **Core Service**: Main promotion logic and business rules
- **Kafka Integration**: Event production and consumption
- **Database Layer**: PostgreSQL for persistent storage
- **Cache Layer**: Redis for high-performance caching
- **Health Monitoring**: System health checks and metrics
- **Promotion Feature**: Specialized discount calculation logic

---

## 🚀 Getting Started

### Prerequisites

- Go 1.23.4 or higher
- Docker and Docker Compose
- PostgreSQL
- golang-migration CLI
- protoc for Golang (optional)
- Redis (optional)
- Kafka (optional)
- Opentelemetry (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/pyre-promotion.git
cd pyre-promotion
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run DB Migration:
```bash
make migrate_up
```


### Project Structure

```
.
├── core/               # Core business logic
├── core-internal/      # Internal core components
├── feature-discount/   # Discount calculation feature
├── feature-health/     # Health monitoring
├── kafka-consume/      # Kafka consumer implementation
├── kafka-produce/      # Kafka producer implementation
├── protos/            # Protocol buffer definitions
└── sqlc/              # SQL code generation
```

---

## 🚢 Deployment

The service can be deployed using:

1. Local Host:
```bash
go run main.go
```

2. Docker:
```bash
docker-compose up --build
```

3. Kubernetes:
```bash
kubectl apply -f deployment.yaml
kubectl apply -f configMap.yaml
```

---

## 💼 Usage
Curls below are example usage of HTTP APIs

### ✏️ Postman API Collection
```bash
https://www.postman.com/dark-eclipse-55522/workspace/pyre-public/collection/20536686-46c3db49-d794-4e99-b60e-6259265e181c?action=share&creator=20536686
```

### ♥️ System Health
```bash
curl --location 'localhost:8000/api/health/v1'
```

### ➕ Insert Promotion
```bash
curl --location 'localhost:8000/api/discount/v1' \
--header 'x-shop-id: shop123' \
--header 'Content-Type: application/json' \
--data '{
  "name": "Holiday Sale 1",
  "promotion_type": "discount",
  "code": "HOLIDAY2024",
  "start_time": "{{start_time}}",
  "end_time": "{{end_time}}",
  "shop_id": "shop123",
  "usage_quantity": 1000,
  "usage_limit_per_user": 5,
  "products": [
    {
      "sku": "product123",
      "name": "Winter Jacket",
      "purchase_limit": 2,
      "product_variants": [
        {
          "sku": "variant123",
          "name": "Winter Jacket - Red",
          "discounted_price": 49.99,
          "discounted_percentage": 20.0,
          "stock_limit": 50,
          "is_active": true
        },
        {
          "sku": "variant124",
          "name": "Winter Jacket - Blue",
          "discounted_price": 44.99,
          "discounted_percentage": 25.0,
          "stock_limit": 30,
          "is_active": true
        }
      ]
    },
    {
      "sku": "product124",
      "name": "Wool Scarf",
      "purchase_limit": 3,
      "product_variants": [
        {
          "sku": "variant125",
          "name": "Wool Scarf - Green",
          "discounted_price": 19.99,
          "discounted_percentage": 10.0,
          "stock_limit": 100,
          "is_active": true
        }
      ]
    }
  ]
}'
```

### 🔍 Detail Promotion
```bash
curl --location 'localhost:8000/api/discount/v1/2dce604e-6d30-4c80-b449-17c67c75dc58' \
--header 'x-shop-id: shop123' \
--data ''
```


### 🔧 Update Promotion
```bash
curl --location --request PUT 'localhost:8000/api/discount/v1/2dce604e-6d30-4c80-b449-17c67c75dc58' \
--header 'x-shop-id: shop123' \
--header 'Content-Type: application/json' \
--data '
{
        "id": 1,
        "created_at": "2024-09-10T03:06:44.779032Z",
        "updated_at": "2024-09-10T03:06:44.779032Z",
        "deleted_at": null,
        "uuid": "2dce604e-6d30-4c80-b449-17c67c75dc58",
        "name": "Holiday Sale 1",
        "promotion_type": "discount",
        "code": "HOLIDAY2024",
        "start_time": "2024-09-09T19:06:44Z",
        "end_time": "2024-09-12T19:06:44Z",
        "shop_id": "shop1234",
        "usage_quantity": 100,
        "usage_limit_per_user": 5,
        "products": [
            {
                "id": 1,
                "created_at": null,
                "updated_at": null,
                "deleted_at": null,
                "uuid": "70f4c414-50df-46b9-a760-8298d7b092c2",
                "promotion_id": "2dce604e-6d30-4c80-b449-17c67c75dc58",
                "sku": "product123",
                "name": "Winter Jacket v2",
                "purchase_limit": 4,
                "product_variants": [
                    {
                        "id": 1,
                        "created_at": null,
                        "updated_at": null,
                        "deleted_at": null,
                        "uuid": "bee08710-6a53-4ef2-956c-2c0e787214b7",
                        "promotion_id": "2dce604e-6d30-4c80-b449-17c67c75dc58",
                        "product_id": "70f4c414-50df-46b9-a760-8298d7b092c2",
                        "sku": "variant123",
                        "name": "Winter Jacket - Red v2",
                        "discounted_price": 49.99,
                        "discounted_percentage": 20,
                        "stock_limit": 500,
                        "is_active": true
                    },
                    {
                        "id": 2,
                        "created_at": null,
                        "updated_at": null,
                        "deleted_at": null,
                        "uuid": "2cc4f257-f6be-4076-968a-4c3808749764",
                        "promotion_id": "2dce604e-6d30-4c80-b449-17c67c75dc58",
                        "product_id": "70f4c414-50df-46b9-a760-8298d7b092c2",
                        "sku": "variant124",
                        "name": "Winter Jacket - Blue v2",
                        "discounted_price": 44.99,
                        "discounted_percentage": 25,
                        "stock_limit": 300,
                        "is_active": true
                    }
                ]
            },
            {
                "id": 2,
                "created_at": null,
                "updated_at": null,
                "deleted_at": null,
                "uuid": "6387f1eb-d399-4e5e-a570-3955a92c89fe",
                "promotion_id": "2dce604e-6d30-4c80-b449-17c67c75dc58",
                "sku": "product124",
                "name": "Wool Scarf v2",
                "purchase_limit": 3,
                "product_variants": [
                    {
                        "id": 3,
                        "created_at": null,
                        "updated_at": null,
                        "deleted_at": null,
                        "uuid": "0a8b8568-a29d-482e-a144-799c81cf9403",
                        "promotion_id": "2dce604e-6d30-4c80-b449-17c67c75dc58",
                        "product_id": "6387f1eb-d399-4e5e-a570-3955a92c89fe",
                        "sku": "variant125",
                        "name": "Wool Scarf - Green v2",
                        "discounted_price": 19.99,
                        "discounted_percentage": 10,
                        "stock_limit": 100,
                        "is_active": true
                    }
                ]
            }
        ]
    }'
```


### 📜 List Promotions
```bash
curl --location 'localhost:8000/api/discount/v1?cursor=5&size=10&sort=updated_at%20ASC' \
--header 'x-shop-id: shop123' \
--data ''
```

### ❌ Delete Promotions
```bash
curl --location --request DELETE 'localhost:8000/api/discount/v1/2dce604e-6d30-4c80-b449-17c67c75dc58' \
--header 'x-shop-id: shop123' \
--data ''
```

---
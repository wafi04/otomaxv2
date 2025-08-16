# Struktur Folder Golang Monolitik - Web Topup

```
topup-backend/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point aplikasi
├── internal/
│   ├── config/
│   │   ├── config.go              # Konfigurasi aplikasi
│   │   └── database.go            # Konfigurasi database
│   ├── domain/
│   │   ├── entities/              # Business entities
│   │   │   ├── user.go
│   │   │   ├── product.go
│   │   │   ├── transaction.go
│   │   │   ├── payment.go
│   │   │   └── topup.go
│   │   ├── repositories/          # Interface repositories
│   │   │   ├── user_repository.go
│   │   │   ├── product_repository.go
│   │   │   ├── transaction_repository.go
│   │   │   └── payment_repository.go
│   │   └── services/              # Business logic interfaces
│   │       ├── auth_service.go
│   │       ├── topup_service.go
│   │       ├── payment_service.go
│   │       └── notification_service.go
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres/
│   │   │   │   ├── connection.go
│   │   │   │   ├── migrations/
│   │   │   │   │   ├── 001_create_users_table.sql
│   │   │   │   │   ├── 002_create_products_table.sql
│   │   │   │   │   ├── 003_create_transactions_table.sql
│   │   │   │   │   └── 004_create_payments_table.sql
│   │   │   │   └── repositories/
│   │   │   │       ├── user_repository_impl.go
│   │   │   │       ├── product_repository_impl.go
│   │   │   │       ├── transaction_repository_impl.go
│   │   │   │       └── payment_repository_impl.go
│   │   │   └── redis/
│   │   │       ├── connection.go
│   │   │       └── cache_service.go
│   │   ├── external/
│   │   │   ├── payment_gateway/
│   │   │   │   ├── midtrans/
│   │   │   │   │   ├── midtrans_client.go
│   │   │   │   │   └── midtrans_webhook.go
│   │   │   │   ├── xendit/
│   │   │   │   │   ├── xendit_client.go
│   │   │   │   │   └── xendit_webhook.go
│   │   │   │   └── gopay/
│   │   │   │       └── gopay_client.go
│   │   │   ├── providers/
│   │   │   │   ├── telkomsel/
│   │   │   │   │   └── telkomsel_api.go
│   │   │   │   ├── indosat/
│   │   │   │   │   └── indosat_api.go
│   │   │   │   ├── xl/
│   │   │   │   │   └── xl_api.go
│   │   │   │   ├── smartfren/
│   │   │   │   │   └── smartfren_api.go
│   │   │   │   ├── steam/
│   │   │   │   │   └── steam_api.go
│   │   │   │   ├── garena/
│   │   │   │   │   └── garena_api.go
│   │   │   │   ├── codashop/
│   │   │   │   │   └── codashop_api.go
│   │   │   │   └── pln/
│   │   │   │       └── pln_api.go
│   │   │   └── notification/
│   │   │       ├── whatsapp/
│   │   │       │   └── whatsapp_client.go
│   │   │       ├── email/
│   │   │       │   └── smtp_client.go
│   │   │       └── sms/
│   │   │           └── sms_client.go
│   │   └── web/
│   │       ├── middleware/
│   │       │   ├── auth.go
│   │       │   ├── cors.go
│   │       │   ├── rate_limiter.go
│   │       │   ├── logger.go
│   │       │   └── recovery.go
│   │       ├── handlers/
│   │       │   ├── auth_handler.go
│   │       │   ├── user_handler.go
│   │       │   ├── product_handler.go
│   │       │   ├── topup_handler.go
│   │       │   ├── transaction_handler.go
│   │       │   ├── payment_handler.go
│   │       │   └── webhook_handler.go
│   │       ├── dto/
│   │       │   ├── request/
│   │       │   │   ├── auth_request.go
│   │       │   │   ├── topup_request.go
│   │       │   │   ├── payment_request.go
│   │       │   │   └── user_request.go
│   │       │   └── response/
│   │       │       ├── auth_response.go
│   │       │       ├── topup_response.go
│   │       │       ├── payment_response.go
│   │       │       ├── user_response.go
│   │       │       └── common_response.go
│   │       └── routes/
│   │           ├── routes.go
│   │           ├── auth_routes.go
│   │           ├── api_routes.go
│   │           └── webhook_routes.go
│   ├── application/
│   │   ├── services/              # Business logic implementation
│   │   │   ├── auth_service_impl.go
│   │   │   ├── topup_service_impl.go
│   │   │   ├── payment_service_impl.go
│   │   │   ├── transaction_service_impl.go
│   │   │   └── notification_service_impl.go
│   │   └── usecases/              # Use case orchestration
│   │       ├── topup_usecase.go
│   │       ├── payment_usecase.go
│   │       └── transaction_usecase.go
│   └── shared/
│       ├── constants/
│       │   ├── status.go
│       │   ├── payment_method.go
│       │   └── provider.go
│       ├── utils/
│       │   ├── jwt.go
│       │   ├── validator.go
│       │   ├── encryption.go
│       │   ├── response.go
│       │   ├── logger.go
│       │   └── helpers.go
│       └── errors/
│           ├── custom_errors.go
│           └── error_codes.go
├── pkg/                           # Packages yang bisa digunakan di luar aplikasi
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   └── crypto/
│       └── crypto.go
├── web/                           # Static files & templates (jika ada)
│   ├── static/
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── templates/
│       ├── layouts/
│       └── pages/
├── scripts/                       # Scripts untuk deployment, migration, dll
│   ├── migrate.sh
│   ├── seed.sh
│   └── deploy.sh
├── tests/                        # Testing files
│   ├── integration/
│   │   ├── auth_test.go
│   │   ├── topup_test.go
│   │   └── payment_test.go
│   ├── unit/
│   │   ├── services/
│   │   ├── handlers/
│   │   └── repositories/
│   └── mocks/                    # Mock objects
│       ├── mock_repositories.go
│       └── mock_services.go
├── docs/                         # Documentation
│   ├── api/
│   │   ├── swagger.yaml
│   │   └── postman/
│   ├── database/
│   │   └── schema.md
│   └── deployment/
│       └── docker.md
├── deployments/                  # Docker & deployment configs
│   ├── docker/
│   │   ├── Dockerfile
│   │   ├── docker-compose.yml
│   │   └── docker-compose.prod.yml
│   └── k8s/                     # Kubernetes manifests (jika diperlukan)
│       ├── deployment.yaml
│       ├── service.yaml
│       └── configmap.yaml
├── .env.example                 # Environment variables template
├── .env                        # Environment variables (don't commit)
├── .gitignore
├── .dockerignore
├── Makefile                    # Build automation
├── go.mod
├── go.sum
└── README.md
```

## Penjelasan Struktur:

### 1. **cmd/server/**

Entry point aplikasi, berisi main.go

### 2. **internal/**

Kode internal aplikasi yang tidak bisa diakses dari luar

### 3. **domain/**

- **entities**: Model bisnis/domain
- **repositories**: Interface untuk data access
- **services**: Interface untuk business logic

### 4. **infrastructure/**

- **database**: Implementasi repository dan koneksi DB
- **external**: Integrasi dengan service eksternal
- **web**: HTTP handlers, middleware, routes

### 5. **application/**

- **services**: Implementasi business logic
- **usecases**: Orchestration use cases

### 6. **shared/**

Utility dan helper yang digunakan di seluruh aplikasi

### 7. **pkg/**

Package yang bisa digunakan di luar aplikasi

### Fitur Topup yang Didukung:

- **Pulsa**: Telkomsel, Indosat, XL, Smartfren
- **Games**: Steam, Garena, Mobile Legends, PUBG
- **E-wallet**: GoPay, OVO, DANA
- **PLN**: Token listrik
- **Internet**: Paket data operator

### Payment Gateway:

- Midtrans
- Xendit
- GoPay
- Bank Transfer
- Virtual Account

Struktur ini mengikuti Clean Architecture dan Domain Driven Design (DDD) yang memisahkan concern dengan baik dan mudah untuk di-maintain serta di-scale.

# IoT Service

The IoT Service provides device management, telemetry collection, and IoT device orchestration for the BengoBox ecosystem. It integrates with auth-service for SSO authentication and user management.

## Key Features

- Multi-tenant IoT device management with organization-aware access control
- User management with RBAC (Role-Based Access Control) and permissions
- SSO integration with auth-service for centralized authentication
- Real-time telemetry collection and processing via event-driven architecture
- Device provisioning and lifecycle management
- RESTful API with OpenAPI/Swagger documentation

## Tech Stack

- Go 1.24+, PostgreSQL, Redis, NATS
- HTTP transport via `chi` router
- JWT validation via `shared/auth-client`
- Observability: zap logging, Prometheus metrics, OpenTelemetry traces

## Getting Started

```shell
cp config/example.env .env
go mod download
docker compose up -d postgres redis nats
go run ./cmd/api
```

APIs default to `http://localhost:4006`. Configure via `IOT_HTTP_PORT`.

## Project Structure

```
cmd/
  api/         # HTTP entrypoint
internal/
  app/         # Bootstrap and lifecycle
  config/      # Environment configuration loader
  http/        # Chi handlers and routes
  platform/    # Infrastructure adapters (database, cache, events)
  services/    # Domain services (rbac, usersync)
  shared/      # Logger and middleware
```

## User Management & RBAC

The service includes comprehensive user management with:
- User creation and synchronization with auth-service SSO
- Role-Based Access Control (RBAC) with permissions
- Tenant-aware access control
- User role assignment and management

### Default Roles

- **admin**: Full access to all devices and telemetry data
- **member**: Can create and manage devices, read telemetry
- **viewer**: Read-only access to devices and telemetry

### Permissions

- `iot:devices:read` - Read IoT devices
- `iot:devices:write` - Create and update devices
- `iot:devices:delete` - Delete devices
- `iot:devices:manage` - Full device management
- `iot:telemetry:read` - Read telemetry data

## API Documentation

- Swagger UI: `http://localhost:4006/v1/docs/` (when implemented)
- All API endpoints are under `/api/v1/{tenantID}/`

## Environment Variables

All configuration keys prefixed with `IOT_`. See [`config/example.env`](config/example.env) for details.

### Auth Service SSO Integration

- `IOT_AUTH_SERVICE_URL`: Auth service URL (default: `https://auth.codevertex.local:4101`)
- `IOT_AUTH_JWKS_URL`: JWKS endpoint for JWT validation
- `IOT_AUTH_SERVICE_API_KEY`: API key for user sync operations

## Documentation

- [`plan.md`](plan.md) - Service architecture and roadmap
- [`CHANGELOG.md`](CHANGELOG.md) - Version history
- [`docs/`](docs/) - Additional documentation

## Community & Governance

- [`CONTRIBUTING.md`](CONTRIBUTING.md)
- [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md)
- [`SECURITY.md`](SECURITY.md)
- [`SUPPORT.md`](SUPPORT.md)


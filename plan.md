# IoT Service - Implementation Plan

## Executive Summary

**System Purpose**: Unified, multi-tenant IoT device management platform for all BengoBox domains (cafe, POS, logistics, inventory) with real-time telemetry collection and processing. Supports device lifecycle management, protocol adapters, and event-driven automation.

**Key Capabilities**:
- Device registration, provisioning, and configuration
- Real-time telemetry collection and processing
- Device communication (MQTT, HTTP, CoAP)
- Rules engine and automation
- Firmware updates and OTA management
- Device health monitoring

**Entity Ownership**: This service owns all IoT-specific entities: devices, device attributes, device groups, telemetry data, device commands, device rules, and firmware updates. **IoT does NOT own**: users (references auth-service via `user_id`), tenants (references from auth-service), inventory items (references from inventory-service).

---

## Technology Stack

### Core Framework
- **Language**: Go 1.24+
- **Architecture**: Clean/Hexagonal architecture
- **HTTP Router**: chi
- **API Documentation**: OpenAPI-first contracts

### Data & Caching
- **Primary Database**: PostgreSQL 16+ with TimescaleDB extension
- **ORM**: Ent (schema-as-code migrations)
- **Caching**: Redis 7+ for caching, rate limiting, device state
- **Message Broker**: NATS JetStream

### Supporting Libraries
- **MQTT**: Eclipse Paho MQTT client
- **Validation**: Custom validators
- **Logging**: zap (structured logging)
- **Tracing**: OpenTelemetry instrumentation
- **Metrics**: Prometheus

### DevOps & Observability
- **Containerization**: Multi-stage Docker builds
- **Orchestration**: Kubernetes (via centralized devops-k8s)
- **CI/CD**: GitHub Actions → ArgoCD
- **Monitoring**: Prometheus + Grafana, OpenTelemetry
- **APM**: Jaeger distributed tracing

---

## Domain Modules & Features

### 1. Device Management

**IoT-Specific Features**:
- Device registration and provisioning
- Device lifecycle (active, inactive, decommissioned)
- Device metadata and attributes
- Device grouping and hierarchies
- Secure device provisioning with certificates

**Entities Owned**:
- `devices` - IoT device registry
- `device_attributes` - Device metadata and configuration
- `device_groups` - Logical device groupings
- `device_group_members` - Device-group relationships
- `device_provisioning` - Device provisioning and certificate management

**Integration Points**:
- **auth-service**: Tenant/outlet discovery
- **inventory-service**: Device inventory tracking

### 2. User Management & RBAC

**IoT-Specific Features**:
- Service-level RBAC roles (admin, member, viewer)
- Permission definitions
- User role assignment
- Tenant-aware access control

**Entities Owned**:
- `user_roles` - User role assignments
- `roles` - Service-level roles
- `permissions` - Permission definitions
- `role_permissions` - Role-permission mappings

**Integration Points**:
- **auth-service**: User identity sync (references only)

### 3. Telemetry Collection

**IoT-Specific Features**:
- Real-time telemetry ingestion
- Data aggregation and processing
- Time-series data storage
- Telemetry querying and analytics

**Entities Owned**:
- `telemetry_data` - Time-series telemetry data
- `telemetry_aggregates` - Pre-aggregated telemetry

**Integration Points**:
- **inventory-service**: Temperature/humidity data for compliance
- **logistics-service**: Device location tracking
- **notifications-service**: Threshold alerts

### 4. Device Communication

**IoT-Specific Features**:
- Protocol adapters (MQTT, HTTP, CoAP)
- Command and control messaging
- Firmware updates and OTA management
- Device heartbeat and health monitoring

**Entities Owned**:
- `device_commands` - Command and control messages
- `command_history` - Command execution history
- `firmware_updates` - OTA firmware update tracking
- `device_heartbeats` - Device health monitoring

**Integration Points**:
- **MQTT Broker**: Device communication
- **External Providers**: Protocol-specific integrations

### 5. Rules & Automation

**IoT-Specific Features**:
- Device rules engine
- Event-driven automation
- Alerting and notifications
- Integration with other services

**Entities Owned**:
- `device_rules` - Device automation rules
- `rule_executions` - Rule execution history

**Integration Points**:
- **notifications-service**: Alert notifications
- **inventory-service**: Compliance alerts
- **logistics-service**: Location-based triggers

---

## Cross-Cutting Concerns

### Testing
- Go test suites with table-driven tests
- Testcontainers for integration testing
- Pact for contract tests
- Performance testing for telemetry ingestion

### Observability
- Structured logging (zap)
- Tracing via OpenTelemetry
- Metrics exported via Prometheus
- Distributed tracing via Tempo/Jaeger

### Security
- OWASP ASVS baseline
- TLS everywhere
- Secrets via Vault/Parameter Store
- Rate limiting & anomaly detection middleware
- JWT validation via auth-service
- Device certificate management

### Scalability
- Stateless HTTP layer
- Background workers via NATS/Redis streams
- TimescaleDB for time-series data
- Caching strategy for device state

### Data Modelling
- Ent schemas as single source of truth
- Tenant/outlet discovery webhooks
- Outbox pattern for reliable domain events
- TimescaleDB hypertables for telemetry

---

## API & Protocol Strategy

- **REST-first**: Versioned routes (`/v1/{tenantID}/devices`), documented via OpenAPI
- **MQTT**: MQTT broker for device communication
- **WebSocket/SSE**: Real-time telemetry streaming
- **Webhooks**: Device events, external provider callbacks
- **Idempotency**: Keys, correlation IDs, distributed tracing context propagation

---

## Compliance & Risk Controls

- Align with Kenya Data Protection Act: explicit consent flows, user data export/delete endpoints, audit logging
- Device security: Certificate management, secure provisioning
- Disaster recovery playbook, RTO/RPO targets (<1 hour)

---

## Sprint Delivery Plan

See `docs/sprints/` folder for detailed sprint plans:
- Sprint 0: Foundations
- Sprint 1: User Management & RBAC
- Sprint 2: Device Management
- Sprint 3: Telemetry Collection
- Sprint 4: Device Communication
- Sprint 5: Rules & Automation
- Sprint 6: Analytics & Hardening
- Sprint 7: Launch & Handover

---

## Runtime Ports & Environments

- **Local development**: Service runs on port **4106**
- **Cloud deployment**: All backend services listen on **port 4000** for consistency behind ingress controllers

---

## References

- [Integration Guide](docs/integrations.md)
- [Entity Relationship Diagram](docs/erd.md)
- [Superset Integration](docs/superset-integration.md)
- [Sprint Plans](docs/sprints/)

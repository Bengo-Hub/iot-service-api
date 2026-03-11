# IoT Service - Integration Guide

## Overview

This document provides detailed integration information for all external services and systems integrated with the IoT service, including internal BengoBox microservices and external third-party services.

---

## Table of Contents

1. [Internal BengoBox Service Integrations](#internal-bengobox-service-integrations)
2. [External Third-Party Integrations](#external-third-party-integrations)
3. [Integration Patterns](#integration-patterns)
4. [Two-Tier Configuration Management](#two-tier-configuration-management)
5. [Event-Driven Architecture](#event-driven-architecture)
6. [Integration Security](#integration-security)
7. [Error Handling & Resilience](#error-handling--resilience)

---

## Internal BengoBox Service Integrations

### Auth Service

**Integration Type**: OAuth2/OIDC + Events + REST

**Use Cases**:
- User authentication and authorization
- JWT token validation
- User identity synchronization
- Tenant/outlet discovery

**Architecture**:
- Uses `shared/auth-client` v0.1.0 library for JWT validation
- All protected `/v1/{tenantID}` routes require valid Bearer tokens

**Events Consumed**:
- `auth.tenant.created` - Initialize tenant in IoT system
- `auth.tenant.updated` - Update tenant metadata
- `auth.user.created` - Create user role assignment
- `auth.user.updated` - Update user identity fields

### Notifications Service

**Integration Type**: Events (NATS) + REST API

**Use Cases**:
- Device alerts
- Telemetry threshold notifications
- System events
- Device offline alerts

**REST API Usage**:
- `POST /v1/{tenantId}/notifications/messages` - Send notification

**Events Published**:
- `iot.device.offline` - Device went offline
- `iot.alert.triggered` - Alert triggered (temperature, security, etc.)
- `iot.telemetry.threshold_exceeded` - Telemetry threshold exceeded
- `iot.device.provisioned` - Device successfully provisioned

### Inventory Service

**Integration Type**: Events (NATS) + REST API

**Use Cases**:
- Device inventory tracking
- Asset management
- Temperature/humidity data for compliance

**REST API Usage**:
- `POST /v1/{tenant}/inventory/assets` - Register device as asset
- `GET /v1/{tenant}/inventory/assets/{id}` - Get asset details

**Events Published**:
- `iot.device.registered` - Device registered (for asset tracking)
- `iot.temperature.alert` - Temperature threshold breach (for inventory compliance)

**Events Consumed**:
- `inventory.asset.updated` - Asset information updated

### Logistics Service

**Integration Type**: Events (NATS) + REST API

**Use Cases**:
- Device location tracking
- Vehicle telemetry integration
- Fleet device management

**REST API Usage**:
- `POST /v1/{tenant}/devices/{id}/location` - Update device location
- `GET /v1/{tenant}/devices/{id}/telemetry` - Get device telemetry

**Events Published**:
- `iot.device.location.updated` - Device location updated
- `iot.vehicle.telemetry` - Vehicle telemetry data

**Events Consumed**:
- `logistics.vehicle.assigned` - Vehicle assigned to fleet (link device)

---

## External Third-Party Integrations

### MQTT Broker

**Purpose**: Device communication protocol

**Configuration** (Tier 1):
- Broker URL: Stored encrypted
- Username/Password: Stored encrypted
- TLS certificates: Stored encrypted

**Use Cases**:
- Device telemetry ingestion
- Command and control messaging
- Device state synchronization

**Protocol**: MQTT 3.1.1 and MQTT 5.0

### HTTP/HTTPS Protocol Adapter

**Purpose**: HTTP-based device communication

**Configuration** (Tier 1):
- API endpoints: Configured per device type
- Authentication tokens: Stored encrypted

**Use Cases**:
- RESTful device APIs
- Webhook callbacks from devices
- Device configuration updates

### CoAP Protocol Adapter

**Purpose**: Constrained Application Protocol for low-power devices

**Configuration** (Tier 1):
- CoAP server endpoints: Configured
- Security credentials: Stored encrypted

**Use Cases**:
- Low-power sensor communication
- Constrained device networks

### Device Manufacturers (Future)

**Purpose**: Direct integration with device manufacturer APIs

**Configuration** (Tier 1):
- Manufacturer API credentials: Stored encrypted
- API endpoints: Configured per manufacturer

**Use Cases**:
- Device provisioning via manufacturer APIs
- Firmware update distribution
- Device health monitoring

---

## Integration Patterns

### 1. REST API Pattern (Synchronous)

**Use Case**: Device registration, configuration updates, command execution

**Implementation**:
- HTTP client with retry logic
- Circuit breaker pattern
- Request timeout (5 seconds default)
- Idempotency keys for mutations

### 2. Event-Driven Pattern (Asynchronous)

**Use Case**: Telemetry ingestion, device events, alert notifications

**Transport**: NATS JetStream

**Flow**:
1. Service publishes event to NATS
2. Subscriber services consume event
3. Process event and update local state
4. Publish response events if needed

**Reliability**:
- At-least-once delivery
- Event deduplication via event_id
- Retry on failure
- Dead letter queue for failed events

### 3. MQTT Pattern (Device Communication)

**Use Case**: Real-time device telemetry and commands

**Implementation**:
- MQTT broker for pub/sub
- QoS levels for message reliability
- Retained messages for device state
- Last Will and Testament for offline detection

### 4. Webhook Pattern (Callbacks)

**Use Case**: External provider callbacks, device webhooks

**Implementation**:
- Webhook endpoints in IoT service
- Signature verification (HMAC-SHA256)
- Retry logic for failed deliveries
- Idempotency handling

---

## Two-Tier Configuration Management

### Tier 1: Developer/Superuser Configuration

**Visibility**: Only developers and superusers

**Configuration Items**:
- MQTT broker credentials
- Protocol adapter API keys
- Device manufacturer API credentials
- Database credentials
- Encryption keys

**Storage**:
- Encrypted at rest in database (AES-256-GCM)
- K8s secrets for runtime
- Vault for production secrets

### Tier 2: Business User Configuration

**Visibility**: Normal system users (tenant admins)

**Configuration Items**:
- Device group settings
- Telemetry thresholds
- Alert preferences
- Rule configurations

**Storage**:
- Plain text in database (non-sensitive)
- Tenant-specific configuration tables

---

## Event-Driven Architecture

### Event Catalog

#### Outbound Events (Published by IoT Service)

**iot.device.registered**
```json
{
  "event_id": "uuid",
  "event_type": "iot.device.registered",
  "tenant_id": "tenant-uuid",
  "timestamp": "2024-12-05T10:30:00Z",
  "data": {
    "device_id": "device-uuid",
    "device_code": "DEV-001",
    "device_type": "temperature_sensor",
    "location": {...}
  }
}
```

**iot.telemetry.received**
```json
{
  "event_id": "uuid",
  "event_type": "iot.telemetry.received",
  "tenant_id": "tenant-uuid",
  "timestamp": "2024-12-05T10:30:00Z",
  "data": {
    "device_id": "device-uuid",
    "metric_name": "temperature",
    "metric_value": 25.5,
    "unit": "celsius",
    "timestamp": "2024-12-05T10:30:00Z"
  }
}
```

**iot.alert.triggered**
```json
{
  "event_id": "uuid",
  "event_type": "iot.alert.triggered",
  "tenant_id": "tenant-uuid",
  "timestamp": "2024-12-05T10:30:00Z",
  "data": {
    "device_id": "device-uuid",
    "alert_type": "temperature_threshold",
    "severity": "high",
    "message": "Temperature exceeded threshold"
  }
}
```

#### Inbound Events (Consumed by IoT Service)

**inventory.asset.updated**
```json
{
  "event_id": "uuid",
  "event_type": "inventory.asset.updated",
  "tenant_id": "tenant-uuid",
  "timestamp": "2024-12-05T10:30:00Z",
  "data": {
    "asset_id": "asset-uuid",
    "device_id": "device-uuid",
    "status": "active"
  }
}
```

---

## Integration Security

### Authentication

**JWT Tokens**:
- Validated via `shared/auth-client` library
- JWKS from auth-service
- Token claims include tenant_id for scoping

**Device Authentication**:
- Certificate-based authentication
- Device tokens (JWT) for API access
- MQTT client certificates

**API Keys** (Service-to-Service):
- Stored in K8s secrets
- Rotated quarterly

### Authorization

**Tenant Isolation**:
- All requests scoped by tenant_id
- Device access isolated per tenant
- Data isolation enforced at database level

### Secrets Management

**Encryption**:
- Secrets encrypted at rest (AES-256-GCM)
- Decrypted only when used
- Key rotation every 90 days

**Device Certificates**:
- Stored encrypted in database
- Certificate rotation support
- Revocation list management

---

## Error Handling & Resilience

### Retry Policies

**Exponential Backoff**:
- Initial delay: 1 second
- Max delay: 30 seconds
- Max retries: 3

### Circuit Breaker

**Implementation**:
- Opens after 5 consecutive failures
- Half-open after 60 seconds
- Closes on successful request

### Device Offline Handling

**Strategy**:
- Heartbeat monitoring
- Offline detection after missed heartbeats
- Buffered telemetry for offline devices
- Automatic sync when device reconnects

### Monitoring

**Metrics**:
- API call latency (p50, p95, p99)
- API call success/failure rates
- Event publishing success rates
- Device online/offline status
- Telemetry ingestion rate

**Alerts**:
- High failure rate (>5%)
- Service unavailability
- Event delivery failures
- Device offline threshold exceeded

---

## References

- [Auth Service Integration](../auth-service/auth-service/docs/integrations.md)
- [Inventory Service Integration](../inventory-service/inventory-api/docs/integrations.md)
- [Logistics Service Integration](../logistics-service/logistics-api/docs/integrations.md)
- [Notifications Service Integration](../notifications-service/notifications-api/docs/integrations.md)


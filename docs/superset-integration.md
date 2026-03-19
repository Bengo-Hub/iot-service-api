# IoT Service - Apache Superset Integration

## Overview

The IoT service integrates with the centralized Apache Superset instance for BI dashboards, analytics, and reporting. Superset is deployed as a centralized service accessible to all BengoBox services.

---

## Architecture

### Service Configuration

**Environment Variables**:
- `SUPERSET_BASE_URL` - Superset service URL
- `SUPERSET_ADMIN_USERNAME` - Admin username (K8s secret)
- `SUPERSET_ADMIN_PASSWORD` - Admin password (K8s secret)
- `SUPERSET_API_VERSION` - API version (default: v1)

**Authentication**:
- Admin credentials used for backend-to-Superset communication
- User authentication via JWT tokens passed to Superset for SSO
- Guest tokens generated for embedded dashboards

---

## Integration Methods

### 1. REST API Client

Backend uses Go HTTP client configured for Superset REST API calls.

**Base Configuration**:
- Base URL: `SUPERSET_BASE_URL/api/v1`
- Default headers: `Content-Type: application/json`
- Authentication: Bearer token from Superset login endpoint
- Retry policy: Exponential backoff (3 retries)
- Circuit breaker: Opens after 5 consecutive failures

**Key API Endpoints**:

**Authentication**:
- `POST /api/v1/security/login` - Login with admin credentials
- `POST /api/v1/security/refresh` - Refresh access token
- `POST /api/v1/security/guest_token/` - Generate guest token for embedding

**Data Sources**:
- `GET /api/v1/database/` - List all data sources
- `POST /api/v1/database/` - Create new data source
- `PUT /api/v1/database/{id}` - Update data source

**Dashboards**:
- `GET /api/v1/dashboard/` - List all dashboards
- `POST /api/v1/dashboard/` - Create new dashboard
- `GET /api/v1/dashboard/{id}` - Get dashboard details

### 2. Database Direct Connection

Superset connects directly to PostgreSQL database via read-only user for data access.

**Connection Configuration**:
- Database type: PostgreSQL with TimescaleDB extension
- Connection string: Provided to Superset via data source API
- Read-only user: `superset_readonly` (created in PostgreSQL)
- Permissions: SELECT only on all tables, no write access
- SSL: Required for production connections

**Read-Only User Setup**:
- Create `superset_readonly` role in PostgreSQL
- Grant CONNECT on database
- Grant USAGE on schema
- Grant SELECT on all tables
- Set default privileges for future tables

**Connection String** (for Superset):
```
postgresql://superset_readonly:password@postgresql.infra.svc.cluster.local:5432/iot_db?sslmode=require
```

**Data Source Creation**:
- Data source created programmatically on application startup
- Connection tested before marking as active
- Data source updated if connection parameters change

---

## Pre-Built Dashboards

### 1. Device Overview Dashboard

**Charts**:
- Total devices (metric)
- Online/offline devices (pie chart)
- Device status distribution (bar chart)
- Device health trends (line chart)
- Device registration over time (line chart)

**Filters**:
- Date range
- Device type
- Device status

**Data Source**: `devices`, `device_heartbeats` tables

### 2. Telemetry Analytics Dashboard

**Charts**:
- Telemetry ingestion rate (metric)
- Telemetry by metric type (bar chart)
- Telemetry trends (line chart)
- Average values by device (table)
- Telemetry distribution (histogram)

**Filters**:
- Date range
- Device selection
- Metric name

**Data Source**: `telemetry_data`, `telemetry_aggregates` tables

### 3. Device Health Dashboard

**Charts**:
- Device uptime (metric)
- Battery level distribution (histogram)
- Signal strength trends (line chart)
- Offline incidents (table)
- Health score trends (line chart)

**Filters**:
- Date range
- Device selection
- Health status

**Data Source**: `device_heartbeats`, `devices` tables

### 4. Rules & Automation Dashboard

**Charts**:
- Active rules (metric)
- Rule execution count (bar chart)
- Rule success rate (line chart)
- Failed rule executions (table)
- Rule trigger frequency (line chart)

**Filters**:
- Date range
- Rule type
- Execution status

**Data Source**: `device_rules`, `rule_executions` tables

### 5. Device Communication Dashboard

**Charts**:
- Command success rate (metric)
- Commands by type (pie chart)
- Command latency (line chart)
- Firmware update status (bar chart)
- Communication protocol distribution (pie chart)

**Filters**:
- Date range
- Device selection
- Command type

**Data Source**: `device_commands`, `command_history`, `firmware_updates` tables

---

## Implementation Details

### Initialization Process

1. Authenticate with Superset using admin credentials
2. Create/update data source pointing to PostgreSQL with TimescaleDB
3. Create/update dashboards for each module:
   - Device Overview Dashboard
   - Telemetry Analytics Dashboard
   - Device Health Dashboard
   - Rules & Automation Dashboard
   - Device Communication Dashboard
4. Log warnings for dashboard creation failures (non-blocking)

### Dashboard Bootstrap

**Backend Endpoint**: `GET /api/v1/dashboards/{module}/embed`

**Process**:
1. Extract tenant ID from context
2. Get dashboard ID for module from Superset
3. Generate guest token with Row-Level Security (RLS) clause filtering by tenant_id
4. Construct embed URL with dashboard ID and guest token
5. Return embed URL with expiration time (5 minutes)

### Row-Level Security (RLS)

**Implementation**:
- Guest tokens include RLS clauses
- RLS filters data by `tenant_id`
- Each tenant sees only their data

---

## Error Handling

### Retry Logic

**Retry Policy**:
- Maximum 3 retry attempts
- Exponential backoff (1s, 2s, 4s delays)
- Retry on 5xx errors or network failures
- Return response on success or after max retries

### Circuit Breaker

**Implementation**:
- Opens after 5 consecutive failures
- Half-open after 60 seconds
- Closes on successful request

### Fallback Strategies

**Superset Unavailable**:
- Return cached dashboard URLs (if available)
- Show static dashboard images
- Log error for monitoring
- Alert operations team

---

## Monitoring

### Metrics

**Integration-Specific Metrics**:
- Superset API call latency (p50, p95, p99)
- Dashboard creation/update success rates
- Guest token generation latency
- Data source connection health

**Prometheus Metrics**:
- `superset_api_call_duration_seconds` - Histogram of API call durations (labeled by endpoint, status)
- `superset_dashboard_views_total` - Counter of dashboard views (labeled by dashboard, tenant)

### Alerts

**Alert Conditions**:
- Superset service unavailability
- High API call failure rate (>5%)
- Dashboard creation failures
- Data source connection failures

---

## Security Considerations

### Authentication & Authorization

- Admin credentials stored in K8s secrets
- Guest tokens expire after 5 minutes
- RLS ensures tenant data isolation
- JWT tokens validated for SSO

### Data Privacy

- Read-only database user
- RLS filters enforce tenant isolation
- Sensitive device data masked in logs
- PII data excluded from dashboards (if applicable)

---

## References

- [Apache Superset REST API Documentation](https://superset.apache.org/docs/api)
- [Superset Deployment Guide](../../devops-k8s/docs/superset-deployment.md)
- [Ordering-Backend Superset Integration](../../../ordering-service/ordering-backend/docs/superset-integration.md)


# IoT Service – Entity Relationship Overview

The IoT Service provides device management, telemetry collection, and IoT device orchestration.  
Schemas are defined via Ent and designed for strict multi-tenancy.

> **Conventions**
> - UUID primary keys.
> - `tenant_id` on every table unless noted.
> - Timestamps use `TIMESTAMPTZ`.
> - Geospatial data stored using PostGIS types where applicable.

---

## User Management & RBAC

| Table | Key Columns | Description |
|-------|-------------|-------------|
| `user_roles` | `id`, `tenant_id`, `user_id`, `role_code`, `assigned_at`, `assigned_by`, `metadata` | User role assignments. References `user_id` from `auth-service` (never stores user accounts). |
| `roles` | `code (PK)`, `name`, `description`, `is_system`, `permissions_json`, `created_at`, `updated_at` | Service-level roles (admin, member, viewer). |
| `permissions` | `id`, `tenant_id`, `name`, `module`, `action`, `resource`, `description`, `created_at` | Permission definitions. |
| `role_permissions` | `id`, `role_code`, `permission_id`, `assigned_at` | Role-permission mappings. |

## Device Management

| Table | Key Columns | Description |
|-------|-------------|-------------|
| `devices` | `id`, `tenant_id`, `device_code`, `device_type`, `name`, `description`, `status`, `location_json`, `geo_point`, `metadata`, `registered_at`, `last_seen_at`, `created_at`, `updated_at` | IoT device registry. |
| `device_attributes` | `id`, `device_id`, `attribute_key`, `attribute_value`, `value_type`, `updated_at` | Device metadata and configuration. |
| `device_groups` | `id`, `tenant_id`, `name`, `description`, `group_type`, `metadata`, `created_at`, `updated_at` | Logical device groupings. |
| `device_group_members` | `id`, `device_group_id`, `device_id`, `added_at` | Device-group relationships. |
| `device_provisioning` | `id`, `tenant_id`, `device_id`, `provisioning_status`, `certificate_url`, `provisioned_at`, `metadata` | Device provisioning and certificate management. |

## Telemetry Collection

| Table | Key Columns | Description |
|-------|-------------|-------------|
| `telemetry_data` | `id`, `tenant_id`, `device_id`, `metric_name`, `metric_value`, `unit`, `timestamp`, `geo_point`, `metadata`, `created_at` | Time-series telemetry data. |
| `telemetry_aggregates` | `id`, `tenant_id`, `device_id`, `metric_name`, `aggregation_type`, `time_bucket`, `value`, `count`, `min_value`, `max_value`, `avg_value`, `created_at` | Pre-aggregated telemetry (hourly/daily). |

## Device Communication

| Table | Key Columns | Description |
|-------|-------------|-------------|
| `device_commands` | `id`, `tenant_id`, `device_id`, `command_type`, `payload_json`, `status`, `sent_at`, `acknowledged_at`, `response_json`, `metadata` | Command and control messages. |
| `command_history` | `id`, `command_id`, `event_type`, `payload`, `occurred_at` | Command execution history. |
| `firmware_updates` | `id`, `tenant_id`, `device_id`, `firmware_version`, `update_status`, `scheduled_at`, `started_at`, `completed_at`, `error_message`, `metadata` | OTA firmware update tracking. |
| `device_heartbeats` | `id`, `device_id`, `heartbeat_at`, `status`, `battery_level`, `signal_strength`, `metadata` | Device health monitoring. |

## Rules & Automation

| Table | Key Columns | Description |
|-------|-------------|-------------|
| `device_rules` | `id`, `tenant_id`, `name`, `rule_type`, `trigger_conditions_json`, `action_json`, `is_active`, `created_at`, `updated_at` | Device automation rules. |
| `rule_executions` | `id`, `rule_id`, `device_id`, `triggered_at`, `execution_status`, `result_json`, `error_message` | Rule execution history. |

## Integrations & Eventing

| Table | Key Columns | Description |
|-------|-------------|-------------|
| `outbox_events` | `id`, `tenant_id`, `aggregate_type`, `aggregate_id`, `event_type`, `payload`, `status`, `attempts`, `last_attempt_at`, `created_at` | Reliable event dispatch to other services. |
| `tenant_sync_events` | `id`, `tenant_id`, `tenant_slug`, `source_service`, `payload`, `synced_at`, `status` | Records inbound tenant discovery callbacks from auth-service. |

## Relationships & External Services

**Entity Ownership**: This service owns IoT device entities. It references (but does not own) entities from other services:
- **Users**: References `user_id` from `auth-service` (never stores user accounts)
- **Tenants**: References `tenant_id`, `tenant_slug` from `auth-service` registry

- `devices` link to `device_attributes`, `device_groups`, `telemetry_data`, and `device_commands`.
- `telemetry_data` feeds into `telemetry_aggregates` for analytics.
- `device_rules` trigger actions based on `telemetry_data` conditions.
- `outbox_events` emit domain events to notifications service (`device.offline`, `telemetry.threshold_exceeded`).
- Tenant/outlet discovery callbacks from auth-service ensure this service can provision devices on demand.

## Seed & Reference Data

- Default roles: `admin`, `member`, `viewer`.
- Default permissions: `iot:devices:read`, `iot:devices:write`, `iot:devices:delete`, `iot:devices:manage`, `iot:telemetry:read`.
- Sample device types and configurations.

---

Update this ERD whenever Ent schemas change. Run `go generate ./internal/ent` and refresh downstream documentation to keep integrations aligned.


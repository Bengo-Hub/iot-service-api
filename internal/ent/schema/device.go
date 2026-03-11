package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Device holds the schema definition for the Device entity.
type Device struct {
	ent.Schema
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("tenant_id", uuid.UUID{}),
		field.String("device_code").
			NotEmpty().
			Unique(),
		field.String("device_type").
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.String("status").
			Default("active"),
		field.JSON("location", map[string]any{}).
			Optional(),
		field.String("geo_point").
			Optional(),
		field.JSON("metadata", map[string]any{}).
			Optional(),
		field.Time("registered_at").
			Default(time.Now),
		field.Time("last_seen_at").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Device.
func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("attributes", DeviceAttribute.Type),
		edge.To("group_members", DeviceGroupMember.Type),
		edge.To("telemetry_data", TelemetryData.Type),
		edge.To("commands", DeviceCommand.Type),
		edge.To("heartbeats", DeviceHeartbeat.Type),
		edge.To("firmware_updates", FirmwareUpdate.Type),
	}
}


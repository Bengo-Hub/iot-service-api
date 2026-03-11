package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DeviceHeartbeat holds the schema definition for the DeviceHeartbeat entity.
type DeviceHeartbeat struct {
	ent.Schema
}

// Fields of the DeviceHeartbeat.
func (DeviceHeartbeat) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("device_id", uuid.UUID{}),
		field.Time("heartbeat_at").
			Default(time.Now),
		field.String("status").
			Default("online"),
		field.Int("battery_level").
			Optional(),
		field.Int("signal_strength").
			Optional(),
		field.JSON("metadata", map[string]any{}).
			Optional(),
	}
}

// Edges of the DeviceHeartbeat.
func (DeviceHeartbeat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("device", Device.Type).
			Ref("heartbeats").
			Field("device_id").
			Unique().
			Required(),
	}
}


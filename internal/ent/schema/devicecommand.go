package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DeviceCommand holds the schema definition for the DeviceCommand entity.
type DeviceCommand struct {
	ent.Schema
}

// Fields of the DeviceCommand.
func (DeviceCommand) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("tenant_id", uuid.UUID{}),
		field.UUID("device_id", uuid.UUID{}),
		field.String("command_type").
			NotEmpty(),
		field.JSON("payload", map[string]any{}).
			Optional(),
		field.String("status").
			Default("pending"),
		field.Time("sent_at").
			Optional(),
		field.Time("acknowledged_at").
			Optional(),
		field.JSON("response", map[string]any{}).
			Optional(),
		field.JSON("metadata", map[string]any{}).
			Optional(),
	}
}

// Edges of the DeviceCommand.
func (DeviceCommand) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("device", Device.Type).
			Ref("commands").
			Field("device_id").
			Unique().
			Required(),
		edge.To("history", CommandHistory.Type),
	}
}

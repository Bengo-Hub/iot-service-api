package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FirmwareUpdate holds the schema definition for the FirmwareUpdate entity.
type FirmwareUpdate struct {
	ent.Schema
}

// Fields of the FirmwareUpdate.
func (FirmwareUpdate) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("tenant_id", uuid.UUID{}),
		field.UUID("device_id", uuid.UUID{}),
		field.String("firmware_version").
			NotEmpty(),
		field.String("update_status").
			Default("pending"),
		field.Time("scheduled_at").
			Optional(),
		field.Time("started_at").
			Optional(),
		field.Time("completed_at").
			Optional(),
		field.String("error_message").
			Optional(),
		field.JSON("metadata", map[string]any{}).
			Optional(),
	}
}

// Edges of the FirmwareUpdate.
func (FirmwareUpdate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("device", Device.Type).
			Ref("firmware_updates").
			Field("device_id").
			Unique().
			Required(),
	}
}

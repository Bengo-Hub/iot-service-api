package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DeviceGroup holds the schema definition for the DeviceGroup entity.
type DeviceGroup struct {
	ent.Schema
}

// Fields of the DeviceGroup.
func (DeviceGroup) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("tenant_id", uuid.UUID{}),
		field.String("name").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.String("group_type").
			Default("custom"),
		field.JSON("metadata", map[string]any{}).
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the DeviceGroup.
func (DeviceGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("members", DeviceGroupMember.Type),
	}
}


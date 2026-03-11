package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DeviceAttribute holds the schema definition for the DeviceAttribute entity.
type DeviceAttribute struct {
	ent.Schema
}

// Fields of the DeviceAttribute.
func (DeviceAttribute) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("device_id", uuid.UUID{}),
		field.String("attribute_key").
			NotEmpty(),
		field.String("attribute_value").
			NotEmpty(),
		field.String("value_type").
			Default("string"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the DeviceAttribute.
func (DeviceAttribute) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("device", Device.Type).
			Ref("attributes").
			Field("device_id").
			Unique().
			Required(),
	}
}


package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DeviceGroupMember holds the schema definition for the DeviceGroupMember entity.
type DeviceGroupMember struct {
	ent.Schema
}

// Fields of the DeviceGroupMember.
func (DeviceGroupMember) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("device_group_id", uuid.UUID{}),
		field.UUID("device_id", uuid.UUID{}),
		field.Time("added_at").
			Default(time.Now),
	}
}

// Edges of the DeviceGroupMember.
func (DeviceGroupMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("device_group", DeviceGroup.Type).
			Ref("members").
			Field("device_group_id").
			Unique().
			Required(),
		edge.From("device", Device.Type).
			Ref("group_members").
			Field("device_id").
			Unique().
			Required(),
	}
}


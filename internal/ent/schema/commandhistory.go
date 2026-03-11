package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// CommandHistory holds the schema definition for the CommandHistory entity.
type CommandHistory struct {
	ent.Schema
}

// Fields of the CommandHistory.
func (CommandHistory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("command_id", uuid.UUID{}),
		field.String("event_type").
			NotEmpty(),
		field.JSON("payload", map[string]any{}).
			Optional(),
		field.Time("occurred_at").
			Default(time.Now),
	}
}

// Edges of the CommandHistory.
func (CommandHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("command", DeviceCommand.Type).
			Ref("history").
			Field("command_id").
			Unique().
			Required(),
	}
}

